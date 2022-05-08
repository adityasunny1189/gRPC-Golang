package main

import (
	"context"
	"log"
	"net"

	pb "grpc_learning/helloworld/protos"

	"google.golang.org/grpc"
)

const (
	port = ":8080"
)

type helloServer struct {
	pb.UnimplementedHelloWorldServer
}

func main() {
	lst, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listed: %v", err)
	}
	log.Printf("Server running on port %v\n", port)
	s := grpc.NewServer()
	pb.RegisterHelloWorldServer(s, &helloServer{})
	if err := s.Serve(lst); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *helloServer) EchoHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Received: %v", in)
	res := &pb.HelloResponse{
		Msg: "Hello " + in.Name,
	}
	return res, nil
}
