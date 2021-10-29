package main

import (
	"context"
	"log"
	"time"

	pb "SDCC/operations"
	"google.golang.org/grpc"
)

const (
	serverAddress1 = "172.17.0.2:50051"
	serverAddress2 = "172.17.0.3:50051"
)

func main() {
	// Set up a connection to the server.
	conn1, err := grpc.Dial(serverAddress1, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn1.Close()

	conn2, err := grpc.Dial(serverAddress2, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn2.Close()

	c1 := pb.NewOperationsClient(conn1)
	c2 := pb.NewOperationsClient(conn2)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	r1, err := c1.Put(ctx, &pb.KeyValue{Key: []byte("abc"), Value: []byte("defa")})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Put(\"abc\", \"defa\"): %s", r1.GetMsg())

	r2, err := c2.Get(ctx, &pb.Key{Key: []byte("abc")})
	if err != nil {
		log.Fatal("Get", err)
	}
	log.Printf("Get(\"abc\"): %s", r2.GetValue())

	r3, err := c2.Put(ctx, &pb.KeyValue{Key: []byte("abc"), Value: []byte("defallo")})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Put(\"abc\", \"defallo\"): %s", r3.GetMsg())

	time.Sleep(1000)

	r4, err := c2.Get(ctx, &pb.Key{Key: []byte("abc")})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Get(\"abc\"): %s", r4.GetValue())

	r1, err = c1.Append(ctx, &pb.KeyValue{Key: []byte("abc"), Value: []byte("ghi")})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Append(\"abc\", \"ghi\"): %s", r1.GetMsg())

	r2, err = c2.Get(ctx, &pb.Key{Key: []byte("abc")})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Get(\"abc\"): %s", r2.GetValue())

	r1, err = c2.Del(ctx, &pb.Key{Key: []byte("abc")})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Del(\"abc\"): %s", r1.GetMsg())
	r2, err = c1.Get(ctx, &pb.Key{Key: []byte("abc")})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Get(\"abc\"): %s", r2.GetValue())
}
