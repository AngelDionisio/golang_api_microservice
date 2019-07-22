package main

import (
	"context"
	"fmt"
	"io"
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

	// doUnarySumOperation(c)

	// doServerStreaming(c)

	doClientStreaming(c)

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

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	log.Printf("Starting a PrimeDecomposition Server Streaming RPC... %v\n", c)
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 12390392840,
	}

	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecomposition RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while receiving stream: %v", err)
		}
		log.Printf("Response from PrimeNumberDecomposition: %v\n", res.GetPrimeFactor())
	}

}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	log.Printf("Starting a ComputeAverage Client Streaming RPC")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while opening stream: %v", err)
	}

	numbers := []int32{3, 5, 9, 54, 23}

	for _, number := range numbers {
		fmt.Printf("Sending number: %v\n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response RPC: %v", err)
	}

	fmt.Printf("The average is: %v", res.GetAverage())
}
