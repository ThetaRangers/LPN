# syntax=docker/dockerfile:1
FROM golang
#FROM alpine


EXPOSE 50051
EXPOSE 42424

WORKDIR /app

#COPY server server
COPY config.json config.json
COPY .aws /root/.aws

COPY cloud ./cloud
COPY data ./data
COPY database ./database
COPY dht ./dht
COPY metadata ./metadata
COPY migration ./migration
COPY operations ./operations
COPY replication ./replication
COPY registerServer ./registerServer
COPY utils ./utils
COPY server.go ./server.go
COPY go.mod go.mod
COPY go.sum go.sum

RUN --mount=type=cache,target=/root/.cache/go-build go build server.go
#RUN go build server.go

CMD ["./server"]