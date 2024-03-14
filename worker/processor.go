package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	db "sgithub.com/techschool/simplebank/db/sqlc"
	"sgithub.com/techschool/simplebank/email"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	Shutdown()
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}
type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
	mailer email.EmailSender
}

func NewRedisTaskProcessor(redisOpts asynq.RedisClientOpt, store db.Store, mail email.EmailSender) TaskProcessor {
	logger := NewLogger()
	redis.SetLogger(logger)
	taskServer := asynq.NewServer(redisOpts, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().Err(err).Str("Type", task.Type()).Bytes("Payload", task.Payload()).Msg("failed to process task sending email")
		}),
		Logger: logger,
	})
	return &RedisTaskProcessor{
		server: taskServer,
		store:  store,
		mailer: mail,
	}
}
func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	return processor.server.Start(mux)
}
func (processor *RedisTaskProcessor) Shutdown() {
	processor.server.Shutdown()
}
