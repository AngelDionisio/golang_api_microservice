package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/simple/golang_api_microservice/calculator/calculatorpb"
	"google.golang.org/grpc"
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
