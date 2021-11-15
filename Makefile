.PHONY: proto clean executable all

all: proto executable
docker: proto executable docker_image

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	operations/operations.proto

executable:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build server.go

docker_image:
	DOCKER_BUILDKIT=1 docker build -t app .

clean:
	rm server
