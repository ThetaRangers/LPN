.PHONY: proto clean executable all

all: proto
docker: proto

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	operations.proto
