package main

import (
	"context"
	"fmt"
	"log"

	"github.com/simple/golang_api_microservice/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Calculator client")

	// define client connection (with insecure as this is not production code)
	cc, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("client could not connect to server: %v", err)
	}
	// defer close client connection
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	doUnarySumOperation(c)

}

// doUnarySumOperation makes a sumRequest to a provided CalculatorServiceClient
func doUnarySumOperation(c calculatorpb.CalculatorServiceClient) {
	log.Printf("Starting a Sumunary RPC... %v\n", c)
	req := &calculatorpb.SumRequest{
		FirstNumber:  5,
		SecondNumber: 40,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}

	log.Printf("Response from Sum: %v\n", res.SumResult)
}
