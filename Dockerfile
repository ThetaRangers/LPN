# syntax=docker/dockerfile:1
FROM alpine

EXPOSE 50051
COPY server server
COPY config.json config.json
COPY .aws /root/.aws
CMD ["./server"]