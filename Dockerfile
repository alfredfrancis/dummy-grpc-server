# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30.0
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
RUN apk add --no-cache protobuf-dev

# Copy proto files
COPY api/ ./api/

# Generate code from proto files
RUN mkdir -p pb
RUN protoc --go_out=pb --go-grpc_out=pb \
    --go_opt=module=github.com/alfredfrancis/dummy-grpc-server/pb \
    --go-grpc_opt=module=github.com/alfredfrancis/dummy-grpc-server/pb \
    api/proto/dummydata.proto

# Copy server code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/grpc-server ./cmd/server

# Final stage
FROM alpine:3.18

# Add CA certificates
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /go/bin/grpc-server /usr/local/bin/grpc-server

# Expose gRPC port
EXPOSE 50051

# Command to run
ENTRYPOINT ["/usr/local/bin/grpc-server"]