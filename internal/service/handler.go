package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	pb "github.com/alfredfrancis/dummy-grpc-server/pb"
)

type DummyDataServer struct {
	pb.UnimplementedDummyDataServiceServer
}

func NewDummyDataServer() *DummyDataServer {
	return &DummyDataServer{}
}

// GetDummyData implements the RPC method
func (s *DummyDataServer) GetDummyData(ctx context.Context, req *pb.DummyDataRequest) (*pb.DummyDataResponse, error) {
	log.Printf("Received GetDummyData request with ID: %v", req.RequestId)
	return generateDummyData(req.RequestId), nil
}

// StreamDummyData implements the streaming RPC method
func (s *DummyDataServer) StreamDummyData(req *pb.DummyDataRequest, stream pb.DummyDataService_StreamDummyDataServer) error {
	log.Printf("Received StreamDummyData request with ID: %v", req.RequestId)

	for i := 0; i < 5; i++ {
		data := generateDummyData(fmt.Sprintf("%s-%d", req.RequestId, i))
		if err := stream.Send(data); err != nil {
			return err
		}
		time.Sleep(time.Second)
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
