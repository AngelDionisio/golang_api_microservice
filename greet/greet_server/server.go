package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/simple/golang_api_microservice/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

// Greet implements the `GreetServiceServer` interface in greet.pb.go (greet protobuff)
// this defines our server, takes context, a pointer to GreetRequest and returns
// a pointer to GreetResponse and an error
func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)
	// we can get GetGreeting because request (req) is a greetRequest and it contains a greeting
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName

	// GreetResponse is a struct that returns key `Result` of type string
	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func main() {
	fmt.Println("Hello from server!")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
