package main

import (
	"context"
	"fmt"
	"log"

	"github.com/simple/golang_api_microservice/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	// close connection when closing
	defer cc.Close()

	// create client
	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("Created client: %f", c)

	doUnary(c)

}

// doUnary executes a GreetRequest and logs results
func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do Unary RPC...")
	// GreetRequest's struct has a pointer to a Greeting, which contains a FirstName and LastName
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Angel",
			LastName:  "Dionisio",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", res.Result)
	// c.Greet(context.Background(), in *GreetRequest, opts ...grpc.CallOption) (*GreetResponse, error)
}
