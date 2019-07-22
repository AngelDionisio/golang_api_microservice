package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/simple/golang_api_microservice/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	// doClientStreaming(c)

	doErrorUnary(c)

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

// doServerStreaming client sends a PrimeNumberDecompositionRequest
// and upon successful response, reads the stream being sent from the server
// until there are no more messages to receive
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

// doClientStreaming client has the ability to send multiple messages to server
// when it finishes sending stream of messages, it closes and receives from server
// and finally it prints the message
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

// doErrorUnary client sends two calls to SquareRoot
// one is a good call, the other shows how errors are displayed
func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	log.Printf("Starting a SquareRoot Unary RPC")

	// correct call
	doErrorCall(c, 10)

	// error call
	doErrorCall(c, -2)
}

// doErrorCall client accepts requests, if integer is a negative number
// it returns an InvalidArgument RPC code.
// this is done by converting the server error, into an status.FromError
// which contain multiple functions for showing error messages, codes
func doErrorCall(c calculatorpb.CalculatorServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: n})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC (user error)
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Printf("Error message from server: %v\n", respErr.Message())
				return
			}
		} else {
			log.Fatalf("Big error calling SquareRoot: %v", err)
			return
		}
	}

	fmt.Printf("Result of square root of %v: %v\n", n, res.GetNumberRoot())
}
