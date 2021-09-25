package main

import (
	"context"
	"github.com/dgraph-io/badger"
	"log"
	"net"

	db "SDCC/database"
	pb "SDCC/operations"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

var database *badger.DB

type server struct {
	pb.UnimplementedOperationsServer
}

func (s *server) Get(_ context.Context, in *pb.Key) (*pb.Value, error) {
	log.Printf("Received: Get(%v)", in.GetKey())
	return &pb.Value{Value: db.Get(in.GetKey(), *database)}, nil
}

func (s *server) Put(_ context.Context, in *pb.KeyValue) (*pb.Ack, error) {
	log.Printf("Received: Put(%v, %v)", in.GetKey(), in.GetValue())
	db.Put(in.GetKey(), in.GetValue(), *database)
	return &pb.Ack{Msg: "Ok"}, nil
}

func (s *server) Append(_ context.Context, in *pb.KeyValue) (*pb.Ack, error) {
	log.Printf("Received: Append(%v, %v)", in.GetKey(), in.GetValue())
	db.Append(in.GetKey(), in.GetValue(), *database)
	return &pb.Ack{Msg: "Ok"}, nil
}

func (s *server) Del(_ context.Context, in *pb.Key) (*pb.Ack, error) {
	log.Printf("Received: Del(%v)", in.GetKey())
	db.Del(in.GetKey(), *database)
	return &pb.Ack{Msg: "Ok"}, nil
}

func main() {
	//Open badger DB
	var err error
	database, err = badger.Open(badger.DefaultOptions("badgerDB"))
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *badger.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(database)

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