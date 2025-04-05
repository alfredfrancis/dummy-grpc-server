.PHONY: all proto build test docker-build docker-run

all: proto build test

proto:
	mkdir -p pb
	protoc --go_out=pb --go-grpc_out=pb \
	--go_opt=module=github.com/alfredfrancis/dummy-grpc-server/pb \
	--go-grpc_opt=module=github.com/alfredfrancis/dummy-grpc-server/pb \
	api/proto/dummydata.proto 

build:
	go build -o bin/server cmd/server/main.go

test:
	go test -v ./...

docker-build:
	docker build -t dummy-grpc-server .

docker-run:
	docker rm -f dummy-grpc-server || true
	docker run -d  -p 50051:50051 --name dummy-grpc-server dummy-grpc-server
