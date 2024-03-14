package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	_ "sgithub.com/techschool/simplebank/doc/statik"
	"sgithub.com/techschool/simplebank/email"
	"sgithub.com/techschool/simplebank/worker"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	db "sgithub.com/techschool/simplebank/db/sqlc"
	"sgithub.com/techschool/simplebank/gapi"
	"sgithub.com/techschool/simplebank/pb"
	"sgithub.com/techschool/simplebank/util"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	config, err := util.LoadConfigDB_Server(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	conn, err := sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("can not connect to the DB")
	}
	store := db.NewStore(conn)
	log.Info().Msgf("config.MigrationUrl : %s", config.MigrationUrl)
	log.Info().Msgf("config.DbSource : %s", config.DbSource)
	runDBMigration(config.MigrationUrl, config.DbSource)
	redisOptions := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOptions)
	waitGroup, ctx := errgroup.WithContext(ctx)
	runTaskProcessor(ctx, waitGroup, redisOptions, store, config)
	runGatewayServer(ctx, waitGroup, config, store, taskDistributor)
	runGrpcServer(ctx, waitGroup, config, store, taskDistributor)
	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("errors from the waiting group")
	}
	defer stop()
	//defer conn.Close()
	abc()
}

func runGrpcServer(ctx context.Context, waitGroup *errgroup.Group, config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("can not initiate")
	}
	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimplebankServer(grpcServer, server)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msgf("can not listen to the port : %s", config.GrpcServerAddress)
	}
	waitGroup.Go(func() error {
		log.Info().Msgf("start gRPC server at %s", listener.Addr().String())

		err = grpcServer.Serve(listener)
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			log.Error().Err(err).Msg("gRPC server failed to serve")
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		abc := <-ctx.Done()
		log.Info().Msgf("Value in Context %v", abc)
		log.Info().Msg("graceful shutdown gRPC server")

		grpcServer.GracefulStop()
		log.Info().Msg("gRPC server is stopped")

		return nil
	})
}

func runGatewayServer(ctx context.Context, waitGroup *errgroup.Group, config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	// Configure custom header matcher to forward headers
	headerMatcher := func(headerName string) (mdName string, ok bool) {
		return strings.ToLower(headerName), true
	}

	// Create options for gRPC-Gateway with custom header matcher
	grpcMux := runtime.NewServeMux(jsonOption, runtime.WithIncomingHeaderMatcher(headerMatcher))
	err = pb.RegisterSimplebankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server")
	}
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create statik fs")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)
	httpServer := &http.Server{
		Handler: gapi.HttpLogger(mux),
		Addr:    config.HttpServerAddress,
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("start HTTP gateway server at %s", httpServer.Addr)
		err = httpServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Error().Err(err).Msg("HTTP gateway server failed to serve")
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		abc := <-ctx.Done()
		log.Info().Msgf("Value in Context %v", abc)
		log.Info().Msg("graceful shutdown HTTP gateway server")

		err := httpServer.Shutdown(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("failed to shutdown HTTP gateway server")
			return err
		}

		log.Info().Msg("HTTP gateway server is stopped")
		return nil
	})
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated successfully")
}
func runTaskProcessor(ctx context.Context, waitGroup *errgroup.Group, redisOpt asynq.RedisClientOpt, store db.Store, config util.Config) {
	sender := email.CreateNewEmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, sender)
	log.Info().Msg("starting task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}

	waitGroup.Go(func() error {
		abc := <-ctx.Done()
		log.Info().Msgf("Value in Context %v", abc)
		log.Info().Msg("graceful shutdown task processor")

		taskProcessor.Shutdown()
		log.Info().Msg("task processor is stopped")

		return nil
	})
}

func abc() <-chan int {
	var ch chan int
	select {
	case <-ch: // Receive from channel
		fmt.Println("Received from channel")
	case ch <- 42: // Send to channel
		fmt.Println("Sent to channel")
	}
	return ch
}
