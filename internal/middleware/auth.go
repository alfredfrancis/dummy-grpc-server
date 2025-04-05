package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	AuthToken  = "dummy-secret-token"
	AuthHeader = "authorization"
)

// TokenAuthInterceptor provides a gRPC interceptor for token authentication
func TokenAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if err := validateAuth(ctx); err != nil {
		return nil, err
	}
	return handler(ctx, req)
}

// StreamAuthInterceptor provides token authentication for streaming RPCs
func StreamAuthInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if err := validateAuth(ss.Context()); err != nil {
		return err
	}
	return handler(srv, ss)
}

// Common authorization logic
func validateAuth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	authHeader, ok := md[AuthHeader]
	if !ok || len(authHeader) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	token := authHeader[0]
	if token != AuthToken {
		return status.Errorf(codes.Unauthenticated, "invalid authorization token")
	}
	return nil
}
