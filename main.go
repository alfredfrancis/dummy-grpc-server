package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	pb "github.com/alfredfrancis/dummy-grpc-server/dummydata"
)

const (
	port        = ":50051"
	AUTH_TOKEN  = "dummy-secret-token"
	AUTH_HEADER = "authorization"
)

type server struct {
	pb.UnimplementedDummyDataServiceServer
}

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

	authHeader, ok := md[AUTH_HEADER]
	if !ok || len(authHeader) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	token := authHeader[0]
	if token != AUTH_TOKEN {
		return status.Errorf(codes.Unauthenticated, "invalid authorization token")
	}
	return nil
}

// Generate a dummy data response
func generateDummyData(requestID string) *pb.DummyDataResponse {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := fmt.Sprintf("data-%d", r.Intn(1000))

	// If request ID is provided, use it in the response
	if requestID != "" {
		id = fmt.Sprintf("data-%s-%d", requestID, r.Intn(100))
	}

	now := time.Now()

	return &pb.DummyDataResponse{
		Id:          id,
		Name:        fmt.Sprintf("Dummy Data %d", r.Intn(100)),
		Value:       r.Int31n(1000),
		Description: "This is randomly generated dummy data for testing purposes",
		Tags:        []string{"dummy", "test", "random", fmt.Sprintf("tag%d", r.Intn(10))},
		CreatedAt: &pb.Timestamp{
			Seconds: now.Unix(),
			Nanos:   int32(now.Nanosecond()),
		},
	}
}

// GetDummyData implements the RPC method
func (s *server) GetDummyData(ctx context.Context, req *pb.DummyDataRequest) (*pb.DummyDataResponse, error) {
	log.Printf("Received GetDummyData request with ID: %v", req.RequestId)
	return generateDummyData(req.RequestId), nil
}

// StreamDummyData implements the streaming RPC method
func (s *server) StreamDummyData(req *pb.DummyDataRequest, stream pb.DummyDataService_StreamDummyDataServer) error {
	log.Printf("Received StreamDummyData request with ID: %v", req.RequestId)

	// Send 5 dummy data responses
	for i := 0; i < 5; i++ {
		streamID := fmt.Sprintf("%s-%d", req.RequestId, i)
		if err := stream.Send(generateDummyData(streamID)); err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func main() {
	// Get port from environment variable or use default
	serverPort := os.Getenv("GRPC_PORT")
	if serverPort == "" {
		serverPort = port
	} else {
		serverPort = ":" + serverPort
	}

	// Create listener
	lis, err := net.Listen("tcp", serverPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create gRPC server with auth interceptors
	s := grpc.NewServer(
		grpc.UnaryInterceptor(TokenAuthInterceptor),
		grpc.StreamInterceptor(StreamAuthInterceptor),
	)

	// enable reflection
	reflection.Register(s)

	// Register service
	pb.RegisterDummyDataServiceServer(s, &server{})

	log.Printf("Server listening at %v", lis.Addr())

	// Start serving
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
