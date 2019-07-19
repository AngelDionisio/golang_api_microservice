package main

import (
	"context"
	"fmt"
	"log"

	"github.com/simple/golang_api_microservice/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello from Client!")

	cc, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to: %v", err)
	}
	// close client connection upon leaving main
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	doUnarySum(c)
}

func doUnarySum(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting sum unary operation RPC...")

	req := &calculatorpb.CalculatorRequest{
		Sum: &calculatorpb.Sum{
			Num_1: 10,
			Num_2: 20,
		},
	}

	res, err := c.Add(context.Background(), req)
	if err != nil {
		log.Fatalf("error calling Add RPC: %v", err)
	}

	log.Printf("Response from Add: %v", res.Result)
}
