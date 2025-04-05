package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/alfredfrancis/dummy-grpc-server/internal/middleware"
	"github.com/alfredfrancis/dummy-grpc-server/internal/service"
	pb "github.com/alfredfrancis/dummy-grpc-server/pb"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server with auth interceptors
	s := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.TokenAuthInterceptor),
		grpc.StreamInterceptor(middleware.StreamAuthInterceptor),
	)

	// Register the service implementation
	pb.RegisterDummyDataServiceServer(s, service.NewDummyDataServer())

	// Register reflection service on gRPC server
	reflection.Register(s)

	// Channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	// Start server in a goroutine
	go func() {
		log.Printf("Server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal
	sig := <-sigChan
	log.Printf("Received %v signal, initiating graceful shutdown", sig)

	// Gracefully stop the server
	s.GracefulStop()
	log.Println("Server stopped gracefully")
}
