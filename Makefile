.PHONY: proto clean executable all

all: proto executable
docker: proto docker_image

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	operations/operations.proto

executable:
	GOOS=linux GOARCH=amd64 go build server.go

docker_image:
	DOCKER_BUILDKIT=1 docker build -t app .

clean:
	rm server

clean_volumes:
	docker container prune
	docker volume rm sdcc_volumeapp1 sdcc_volumeapp2 sdcc_volumeapp3 sdcc_volumeapp4 sdcc_volumeapp5 sdcc_volumeapp6 sdcc_volumeapp7 sdcc_volumeapp8 sdcc_volumeapp9 sdcc_volumeapp10