# syntax=docker/dockerfile:1
FROM ubuntu

EXPOSE 50051
COPY server server
#CMD ["./server"]