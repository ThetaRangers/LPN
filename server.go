package main

import (
	"SDCC/utils"
	"context"
	"log"
	"net"

	db "SDCC/database"
	pb "SDCC/operations"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

var database db.Database

type server struct {
	pb.UnimplementedOperationsServer
}

func (s *server) Get(_ context.Context, in *pb.Key) (*pb.Value, error) {
	log.Printf("Received: Get(%v)", in.GetKey())
	return &pb.Value{Value: database.Get(in.GetKey())}, nil
}

func (s *server) Put(_ context.Context, in *pb.KeyValue) (*pb.Ack, error) {
	log.Printf("Received: Put(%v, %v)", in.GetKey(), in.GetValue())
	database.Put(in.GetKey(), in.GetValue())
	return &pb.Ack{Msg: "Ok"}, nil
}

func (s *server) Append(_ context.Context, in *pb.KeyValue) (*pb.Ack, error) {
	log.Printf("Received: Append(%v, %v)", in.GetKey(), in.GetValue())
	database.Append(in.GetKey(), in.GetValue())
	return &pb.Ack{Msg: "Ok"}, nil
}

func (s *server) Del(_ context.Context, in *pb.Key) (*pb.Ack, error) {
	log.Printf("Received: Del(%v)", in.GetKey())
	database.Del(in.GetKey())
	return &pb.Ack{Msg: "Ok"}, nil
}

func init() {
	database = utils.GetConfiguration().Database
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOperationsServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
