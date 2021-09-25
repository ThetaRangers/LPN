package main

import (
	"context"
	"log"
	"time"

	pb "SDCC/operations"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewOperationsClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r1, err := c.Put(ctx, &pb.KeyValue{Key: []byte("abc"), Value: []byte("def")})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Put(\"abc\", \"def\"): %s", r1.GetMsg())
	r2, err := c.Get(ctx, &pb.Key{Key: []byte("abc")})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Get(\"abc\"): %s", r2.GetValue())
	r1, err = c.Append(ctx, &pb.KeyValue{Key: []byte("abc"), Value: []byte("ghi")})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Append(\"abc\", \"ghi\"): %s", r1.GetMsg())
	r2, err = c.Get(ctx, &pb.Key{Key: []byte("abc")})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Get(\"abc\"): %s", r2.GetValue())
	r1, err = c.Del(ctx, &pb.Key{Key: []byte("abc")})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Del(\"abc\"): %s", r1.GetMsg())
	r2, err = c.Get(ctx, &pb.Key{Key: []byte("abc")})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Get(\"abc\"): %s", r2.GetValue())
}
