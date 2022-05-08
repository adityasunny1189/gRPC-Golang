package main

import (
	"context"
	"fmt"
	pb "grpc_learning/helloworld/protos"
	"log"
	"time"

	"google.golang.org/grpc"
)

const (
	address = "localhost:8080"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal("connection failed")
	}
	defer conn.Close()
	client := pb.NewHelloWorldClient(conn)
	EchoHello(client)
}

func EchoHello(client pb.HelloWorldClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.HelloRequest{
		Name: "Aditya Pathak",
	}
	res, err := client.EchoHello(ctx, req)
	if err != nil {
		log.Fatalf("error %v", err)
	}
	fmt.Println("Received: %v", res)
}
