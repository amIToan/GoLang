package gapi

import (
	"context"

	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	db "sgithub.com/techschool/simplebank/db/sqlc"
	"sgithub.com/techschool/simplebank/pb"
	"sgithub.com/techschool/simplebank/validator"
)

// CreateUser implements
func (s *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	violations := validatorVerifyEmailRequest(req)
	if len(violations) > 0 {
		return nil, invalidArgumentError(violations)
	}
	//create email
	user, mail, err := s.store.VerifyEmailTx(ctx, db.UpdateVerifyEmailParams{ID: req.GetId(), SecretCode: req.GetSecretCode()})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "verify mail not implemented")
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to verify: %s", err)
	}
	rsp := pb.VerifyEmailResponse{
		User: convertUser(*user),
		Mail: convertMail(mail),
	}
	return &rsp, nil
}

func validatorVerifyEmailRequest(req *pb.VerifyEmailRequest) []*errdetails.BadRequest_FieldViolation {
	// Get the type of the struct
	var errors []*errdetails.BadRequest_FieldViolation
	if err := validator.ValidateEmailId(req.GetId()); err != nil {
		errors = append(errors, fieldViolation("Id", err))
	}
	if err := validator.ValidateSecretCode(req.GetSecretCode()); err != nil {
		errors = append(errors, fieldViolation("Password", err))
	}
	return errors
}
func convertMail(mail *db.VerifyEmail) *pb.Mail {
	return &pb.Mail{
		Id:        mail.ID,
		Email:     mail.Email,
		IsUsed:    mail.IsUsed,
		CreatedAt: timestamppb.New(mail.CreatedAt),
		ExpiredAt: timestamppb.New(mail.CreatedAt),
	}
}
