package interceptors

import (
	"context"
	"gopasskeeper/internal/lib/jwt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ctxUserID string

const CtxUserID ctxUserID = "ctxUserID"

// AuthInterceptor is a server interceptor for authentication and authorization
type AuthInterceptor struct {
	jwtManager      *jwt.JWTManager
	accessibleRoles map[string][]string
}

// NewAuthInterceptor returns a new auth interceptor
func NewAuthInterceptor(jwtManager *jwt.JWTManager, accessibleRoles map[string][]string) *AuthInterceptor {
	return &AuthInterceptor{jwtManager, accessibleRoles}
}

// Unary returns a server interceptor function to authenticate and authorize unary RPC
func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		userID, err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, CtxUserID, userID)

		return handler(ctx, req)
	}
}

// Stream returns a server interceptor function to authenticate and authorize stream RPC
func (interceptor *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv any,
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		_, err := interceptor.authorize(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) (string, error) {
	accessibleRoles, ok := interceptor.accessibleRoles[method]
	if !ok {
		return "", nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	claims, err := interceptor.jwtManager.Verify(accessToken)
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	for _, role := range accessibleRoles {
		if role == claims.Role {
			return claims.UserID, nil
		}
	}

	return "", status.Error(codes.PermissionDenied, "no permission to access this RPC")
}

func ExtractUID(ctx context.Context) (string, error) {
	uidVal := ctx.Value(CtxUserID)
	uid, ok := uidVal.(string)
	if !ok {
		return "", status.Error(codes.InvalidArgument, "failed extract uid")
	}

	if uid == "" {
		return "", status.Error(codes.DataLoss, "failed extract uid")
	}

	return uid, nil
}
