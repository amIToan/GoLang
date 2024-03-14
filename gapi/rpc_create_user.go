package gapi

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	db "sgithub.com/techschool/simplebank/db/sqlc"
	"sgithub.com/techschool/simplebank/pb"
	"sgithub.com/techschool/simplebank/util"
	"sgithub.com/techschool/simplebank/validator"
	"sgithub.com/techschool/simplebank/worker"
)

// CreateUser implements
func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validatorCreateUserRequest(req)
	if len(violations) > 0 {
		return nil, invalidArgumentError(violations)
	}
	hashedPass, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
	}
	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.GetUsername(),
			HashedPassword: hashedPass,
			FullName:       req.GetFullName(),
			Email:          req.GetEmail(),
		},
		CallbackAfterCreate: func(user db.User) error {
			//todo : use db transaction
			emailTaskPayload := &worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}
			disOpts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.Queue(worker.QueueCritical),
				asynq.ProcessIn(10 * time.Second),
			}
			err = s.taskDistributor.DistributeTaskSendEmail(ctx, emailTaskPayload, disOpts...)
			return err
		},
	}
	log.Info().Msg(">>>>> creating user...")
	time.Sleep(10 * time.Second)
	//create account
	userTxResult, err := s.store.CreateUserTx(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "method CreateUser not implemented")
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}
	rsp := pb.CreateUserResponse{
		User: convertUser(userTxResult.User),
	}
	return &rsp, nil
}

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
func validatorCreateUserRequest(req *pb.CreateUserRequest) []*errdetails.BadRequest_FieldViolation {
	// Get the type of the struct
	var errors []*errdetails.BadRequest_FieldViolation
	if err := validator.ValidateUsername(req.GetUsername()); err != nil {
		errors = append(errors, fieldViolation("username", err))
	}
	if err := validator.ValidatePassword(req.GetPassword()); err != nil {
		errors = append(errors, fieldViolation("Password", err))
	}
	if err := validator.ValidateFullName(req.GetFullName()); err != nil {
		errors = append(errors, fieldViolation("FullName", err))
	}
	if err := validator.ValidateEmail(req.GetEmail()); err != nil {
		errors = append(errors, fieldViolation("Email", err))
	}
	return errors
}
