# syntax=docker/dockerfile:1
FROM golang
#FROM alpine


EXPOSE 50051

WORKDIR /app

#COPY server server
COPY config.json config.json
COPY .aws /root/.aws

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build go build server.go
#RUN go build server.go

CMD ["./server"]