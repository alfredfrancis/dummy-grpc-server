build_proto:
	protoc --proto_path=protos \
	--go_out=dummydata \
	--go-grpc_out=dummydata \
	--go_opt=module=github.com/alfredfrancis/dummy-grpc-server/dummydata \
	--go-grpc_opt=module=github.com/alfredfrancis/dummy-grpc-server/dummydata \
	protos/dummydata.proto

build:
	docker build -t grpc-dummy-server .

run:
	docker run -d -p 50051:50051 grpc-dummy-server

req:
	grpcurl -plaintext -d '{"requestId": "bilfa-red"}' 127.0.0.1:50051 dummydata.DummyDataService.GetDummyData

steaming-req:
	grpcurl -plaintext -d '{"requestId": "alfa-red"}' 127.0.0.1:50051 dummydata.DummyDataService.StreamDummyData