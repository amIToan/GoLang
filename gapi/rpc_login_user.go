package gapi

import (
	"context"
	"database/sql"

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
func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	violations := validatorUserLoginRequest(req)
	if len(violations) > 0 {
		return nil, invalidArgumentError(violations)
	}
	user, err := s.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "Not registered yet!!!")
		}
		return nil, status.Errorf(codes.Internal, "failed to login: %s", err)
	}
	err = util.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "username or password wrong: %s", err)
	}
	accessToken, accessTokenPayload, err := s.tokenMaker.CreateToken(user.Username, util.BankerRole, s.config.ValidDurationTime)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to login: %s", err)
	}
	refreshToken, refreshTokenPayload, err := s.tokenMaker.CreateToken(user.Username, util.BankerRole, s.config.RefreshTokenDurationTime)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to login: %s", err)
	}
	mtdt := s.extractMetadata(ctx)
	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to login: %s", err)
	}
	res := pb.LoginResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiredAt:  timestamppb.New(accessTokenPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: timestamppb.New(refreshTokenPayload.ExpiredAt),
		UserInfo:              convertUser(user),
	}
	return &res, nil
}
func validatorUserLoginRequest(req *pb.LoginRequest) []*errdetails.BadRequest_FieldViolation {
	// Get the type of the struct
	var errors []*errdetails.BadRequest_FieldViolation
	if err := validator.ValidateUsername(req.GetUsername()); err != nil {
		errors = append(errors, fieldViolation("username", err))
	}
	if err := validator.ValidatePassword(req.GetPassword()); err != nil {
		errors = append(errors, fieldViolation("Password", err))
	}
	return errors
}
