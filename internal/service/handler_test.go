package service

import (
	"context"
	"testing"

	pb "github.com/alfredfrancis/dummy-grpc-server/pb"
	"google.golang.org/grpc"
)

type mockStream struct {
	grpc.ServerStream
	results []*pb.DummyDataResponse
}

func (m *mockStream) Send(resp *pb.DummyDataResponse) error {
	m.results = append(m.results, resp)
	return nil
}

func (m *mockStream) Context() context.Context {
	return context.Background()
}

func TestGetDummyData(t *testing.T) {
	server := NewDummyDataServer()
	req := &pb.DummyDataRequest{RequestId: "test-123"}

	resp, err := server.GetDummyData(context.Background(), req)
	if err != nil {
		t.Errorf("GetDummyData failed: %v", err)
	}

	// Verify response fields
	if resp == nil {
		t.Fatal("Expected non-nil response")
	}
	if resp.Id == "" {
		t.Error("Expected non-empty ID")
	}
	if resp.CreatedAt == nil {
		t.Error("Expected non-nil CreatedAt timestamp")
	}
	if len(resp.Tags) == 0 {
		t.Error("Expected non-empty tags")
	}
}

func TestStreamDummyData(t *testing.T) {
	server := NewDummyDataServer()
	req := &pb.DummyDataRequest{RequestId: "test-stream-123"}

	mockStream := &mockStream{}
	err := server.StreamDummyData(req, mockStream)
	if err != nil {
		t.Errorf("StreamDummyData failed: %v", err)
	}

	// Verify stream results
	if len(mockStream.results) != 5 {
		t.Errorf("Expected 5 stream results, got %d", len(mockStream.results))
	}

	// Verify each response in the stream
	for i, resp := range mockStream.results {
		if resp == nil {
			t.Errorf("Stream result %d is nil", i)
			continue
		}
		if resp.Id == "" {
			t.Errorf("Stream result %d has empty ID", i)
		}
	}
}
