package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"sort"

	"github.com/angel/golang_api_microservice/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// this server struct needs to implement the functions of the
// target interface (in this case "CalculatorServiceServer")
// which makes this server also of type CalculatorServiceServer
// we can achieve this by creating the functions defined
// by CalculatorServiceServer with methods
type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Printf("Recieved Sum RPC: %v\n", req)
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber

	sum := firstNumber + secondNumber

	res := &calculatorpb.SumResponse{
		SumResult: sum,
	}

	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	log.Printf("Recieved PrimeNumberDecomposition RPC: %v\n", req)

	number := req.GetNumber()
	divisor := int64(2)

	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			number = number / divisor
		} else {
			divisor++
			fmt.Printf("Divisor has increated to %v\n", divisor)
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	log.Printf("Recieved ComputeAverage RPC\n")

	sum := int32(0)
	count := 0

	for {
		req, err := stream.Recv()
		// when finished getting all requests, then sendAndClose the stream with the ComputeAverageResponse
		if err == io.EOF {
			average := float64(sum) / float64(count)
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})
		}
		if err != nil {
			log.Printf("Error while reading client stream: %v\n", err)
		}

		sum += req.GetNumber()
		count++
	}

}

// MaxIntSlice return the maximum interger in a slice
func MaxIntSlice(v []int) int {
	sort.Ints(v)
	return v[len(v)-1]
}

// MinIntSlice return the minimum interger in a slice
func MinIntSlice(v []int) int {
	sort.Ints(v)
	return v[0]
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("Received FindMaximum RPC call\n")

	currentMax := int32(0)

	// receive messages from client
	for {
		req, err := stream.Recv()
		// when we receive the last of the messages, return the max
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Error while reading FindMaximum client stream: %v\n", err)
			return err
		}

		incomingNumber := req.GetNumber()
		fmt.Printf("Server recieved a new number from stream...: %v\n", incomingNumber)
		// compare incoming value to current max, if its greater, send this new max
		if incomingNumber > currentMax {
			currentMax = incomingNumber
			sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
				Maximum: currentMax,
			})
			if sendErr != nil {
				log.Fatalf("Error while trying to send data to client: %v", sendErr)
				return sendErr
			}
		}
	}

}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Println("Received SquareRoot RPC")

	number := req.GetNumber()

	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("received a negative number: %v", number),
		)
	}

	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	fmt.Println("Calculator Server called")

	// set which port to listen on
	lis, err := net.Listen("tcp", "0.0.0.0:50053")
	if err != nil {
		log.Fatalf("Server could not listen on specified port: %v", err)
	}

	// create new server, that listens on specified port
	s := grpc.NewServer()
	// registering calculator service, this must be invoked before
	// calling Serve
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("server failed to serve: %v", err)
	}
}
