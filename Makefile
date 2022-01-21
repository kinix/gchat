export GCHAT_SERVER ?= 0.0.0.0:50051

protoc:
	protoc gchatpb/gchat.proto --go_out=gchatpb --go-grpc_out=gchatpb

env:

run-server: env
	go run ./server

run-client: env
	go run ./client