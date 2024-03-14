package gapi

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	db "sgithub.com/techschool/simplebank/db/sqlc"
	"sgithub.com/techschool/simplebank/pb"
	"sgithub.com/techschool/simplebank/util"
	"sgithub.com/techschool/simplebank/validator"
)

// implement gRPC login
func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	authPayload, _uuid, err := s.authorization(ctx, []string{util.DepositorRole, util.BankerRole})
	if (err != nil) || (_uuid == uuid.Nil) {
		return nil, unauthenticatedError(err)
	}
	violations := validatorUserUpdateRequest(req)
	if len(violations) > 0 {
		return nil, invalidArgumentError(violations)
	}
	if authPayload.Username != util.BankerRole && authPayload.Username != req.GetUsername() {
		return nil, unauthenticatedError(err)
	}
	oldUser, err := s.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "Not registered yet!!!")
		}
		return nil, status.Errorf(codes.Internal, "failed to update: %s", err)
	}
	var isPassChanged bool
	arg := db.UpdateUserParams{}
	if req.Password != nil {
		err = util.CheckPassword(req.GetPassword(), oldUser.HashedPassword)
		if err != nil {
			isPassChanged = true
		} else {
			isPassChanged = false
		}
		hashedPass, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Unimplemented, "failed to update")
		}
		arg.HashedPassword = sql.NullString{
			String: hashedPass,
			Valid:  true,
		}
	}
	if req.Fullname != nil {
		arg.FullName = sql.NullString{
			String: req.GetFullname(),
			Valid:  true,
		}
	}
	if req.Email != nil {
		arg.Email = sql.NullString{
			String: req.GetEmail(),
			Valid:  true,
		}
	}
	arg.Username = req.Username
	if isPassChanged {
		_, err := s.store.UpdateSession(ctx, db.UpdateSessionParams{
			ID:        _uuid,
			IsBlocked: true,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update session: %s", err)
		}
	}
	updatedUser, err := s.store.UpdateUser(context.Background(), arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update: %s", err)
	}

	res := convertUpdatedUser(updatedUser)
	return res, nil
}
func validatorUserUpdateRequest(req *pb.UpdateUserRequest) []*errdetails.BadRequest_FieldViolation {
	// Get the type of the struct
	var errors []*errdetails.BadRequest_FieldViolation
	if req.Password != nil {
		if err := validator.ValidatePassword(req.GetPassword()); err != nil {
			errors = append(errors, fieldViolation("Password", err))
		}
	}
	if req.Email != nil {
		if err := validator.ValidateEmail(req.GetEmail()); err != nil {
			errors = append(errors, fieldViolation("Email", err))
		}
	}
	if req.Fullname != nil {
		if err := validator.ValidateFullName(req.GetFullname()); err != nil {
			errors = append(errors, fieldViolation("full name", err))
		}
	}
	return errors
}
func convertUpdatedUser(user db.User) *pb.UpdateUserResponse {
	return &pb.UpdateUserResponse{
		User: &pb.User{
			Username:          user.Username,
			FullName:          user.FullName,
			Email:             user.Email,
			PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
			CreatedAt:         timestamppb.New(user.CreatedAt),
		}}
}
