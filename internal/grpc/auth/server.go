package auth

import (
	"context"
	"errors"

	ssov1 "gopasskeeper/internal/grpc/auth/gen/sso"

	"gopasskeeper/internal/services/auth"
	"gopasskeeper/internal/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

// Auth is an interface to describe auth service methods.
type Auth interface {
	Login(
		ctx context.Context,
		login string,
		password string,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		login string,
		password string,
	) (userID string, err error)
}

// Register is a function to bind auth service to a grpc server.
func Register(gRPCServer *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

// Login is an auth server api method to login user with a JWT.
func (s *serverAPI) Login(
	ctx context.Context,
	in *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	if in.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	token, err := s.auth.Login(ctx, in.GetLogin(), in.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid login or password")
		}

		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

// Register is an auth server api method to register user within the service.
func (s *serverAPI) Register(
	ctx context.Context,
	in *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	if in.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	uid, err := s.auth.RegisterNewUser(ctx, in.GetLogin(), in.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &ssov1.RegisterResponse{UserId: uid}, nil
}
