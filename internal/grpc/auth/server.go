package auth

import (
	"auth-microserivce/internal/domain/auth"
	"context"
	"errors"
	v3 "github.com/baneb1ade/auth-protos/gen/go"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var v = validator.New()

const emptyValue = ""

type Service interface {
	Login(ctx context.Context, username, password string) (string, error)
	Register(ctx context.Context, email, username, password string) (string, error)
}

type serverAPI struct {
	v3.UnimplementedAuthServer
	service Service
}

func Register(gRPC *grpc.Server, service Service) {
	v3.RegisterAuthServer(gRPC, &serverAPI{service: service})
}

func (s *serverAPI) Login(ctx context.Context, req *v3.LoginRequest) (*v3.LoginResponse, error) {
	err := validateFields(req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, validateError(err)
	}

	token, err := s.service.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, validateError(err)
	}
	return &v3.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *v3.RegisterRequest) (*v3.RegisterResponse, error) {
	email := req.GetEmail()
	if err := v.Var(email, "required,email"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid email")
	}
	err := validateFields(req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, validateError(err)
	}

	id, err := s.service.Register(ctx, email, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, validateError(err)
	}
	return &v3.RegisterResponse{
		UserId: id,
	}, nil
}

func validateFields(username string, password string) error {
	if password == emptyValue {
		return auth.ErrInvalidPassword
	}
	if username == emptyValue {
		return auth.ErrInvalidUsername
	}
	return nil
}

func validateError(err error) error {
	switch {
	case errors.Is(err, auth.ErrEmailAlreadyExists):
		return status.Error(codes.AlreadyExists, auth.ErrEmailAlreadyExists.Error())

	case errors.Is(err, auth.ErrUsernameAlreadyExists):
		return status.Error(codes.AlreadyExists, auth.ErrUsernameAlreadyExists.Error())

	case errors.Is(err, auth.ErrInvalidUsernameOrPassword):
		return status.Error(codes.Unauthenticated, auth.ErrInvalidUsernameOrPassword.Error())

	case errors.Is(err, auth.ErrInvalidPassword):
		return status.Error(codes.InvalidArgument, auth.ErrInvalidPassword.Error())

	case errors.Is(err, auth.ErrInvalidUsername):
		return status.Error(codes.InvalidArgument, auth.ErrInvalidUsername.Error())

	case errors.Is(err, auth.ErrInvalidEmail):
		return status.Error(codes.InvalidArgument, auth.ErrInvalidEmail.Error())

	default:
		return status.Error(codes.Internal, auth.ErrSmthWentWrong.Error())
	}
}
