package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "SDCC/operations"
	"google.golang.org/grpc"
)

const (
	serverAddress1 = "172.17.0.2:50051"
	serverAddress2 = "172.17.0.4:50051"
)

func ping(address string) time.Duration {

	conn, _ := grpc.Dial(address+":50051", grpc.WithInsecure(), grpc.WithBlock())
	c := pb.NewOperationsClient(conn)
	defer conn.Close()

	start := time.Now()
	c.Ping(context.Background(), &pb.PingMessage{})
	elapsed := time.Since(start)

	return elapsed
}

func main() {
	// Set up a connection to the server.
	fmt.Println("Ping: ", ping("172.17.0.2"))

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
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	input := make([][]byte, 1)

	grr, err := c1.Get(ctx, &pb.Key{Key: []byte("abc")})
	if err != nil {
		log.Fatal("Get", err)
	}
	log.Printf("Get(\"abc\"): %s", grr.GetValue())

	input[0] = []byte("defa")
	r1, err := c1.Put(ctx, &pb.KeyValue{Key: []byte("abc"), Value: input})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Put(\"abc\", \"defa\"): %s", r1.GetMsg())

	time.Sleep(3 * time.Second)
	r2, err := c2.Get(ctx, &pb.Key{Key: []byte("abc")})
	if err != nil {
		log.Fatal("Get", err)
	}
	log.Printf("Get(\"abc\"): %s", r2.GetValue())

	input[0] = []byte("defallo")
	r3, err := c2.Put(ctx, &pb.KeyValue{Key: []byte("abc"), Value: input})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Put(\"abc\", \"defallo\"): %s", r3.GetMsg())
	time.Sleep(3 * time.Second)

	r4, err := c2.Get(ctx, &pb.Key{Key: []byte("abc")})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Get(\"abc\"): %s", r4.GetValue())

	input[0] = []byte("ghi")
	r1, err = c2.Append(ctx, &pb.KeyValue{Key: []byte("abc"), Value: input})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Append(\"abc\", \"ghi\"): %s", r1.GetMsg())
	time.Sleep(3 * time.Second)

	r2, err = c2.Get(ctx, &pb.Key{Key: []byte("abc")})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Get(\"abc\"): %s", r2.GetValue())

	/*
		r1, err = c2.Del(ctx, &pb.Key{Key: []byte("abc")})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Del(\"abc\"): %s", r1.GetMsg())*/
	r2, err = c2.Get(ctx, &pb.Key{Key: []byte("abc")})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Get(\"abc\"): %s", r2.GetValue())

	r2, err = c1.Get(ctx, &pb.Key{Key: []byte("abbacchio")})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Get(\"abbacchio\"): %s", r2.GetValue())
}
