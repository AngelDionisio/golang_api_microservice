package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/simple/golang_api_microservice/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	// doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBiDiStreaming(c)

	doUnaryCallWithDeadline(c, 5*time.Second) // should complete
	doUnaryCallWithDeadline(c, 1*time.Second) // should timeout

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

// doServerStreaming sends a GreetManyTimesRequest and prints the responses from the server
// while there are messages to be sent. Using io.EOF (io.EndOfFile) to break out of infinite for loop
func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a service streaming RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Angel",
			LastName:  "Dionisio",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Angel",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Luz",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Julio",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Jose",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet RPC: %v", err)
	}

	// we iterate over our slice and send each message independently
	for _, req := range requests {
		fmt.Printf("Sending requests: %v\n", req)
		stream.Send(req)
		time.Sleep(time.Millisecond * 1000)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet RPC: %v", err)
	}

	fmt.Printf("LongGreet response: %v\n", res)

}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC...")

	// create stream by invoking client, which returns a client and an error
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while creating stream: %v", err)
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Angel",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Luz",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Julio",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Jose",
			},
		},
	}

	doneChannel := make(chan struct{})
	// send a bunch of messages to client (as goRoutines)
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(time.Millisecond * 1000) // sleep for one second
		}
		stream.CloseSend()
	}()

	// receive a bunch of messages from client (goRoutines)
	// when looking to know the function definition of a client or server
	// look at the "<name>Server interface" on the service.pb.go
	go func() {
		for {
			res, err := stream.Recv()
			// when we finish receiving all messages from the server,
			// then close the doneChannel to signal we are done, and stop blocking
			if err == io.EOF {
				close(doneChannel)
			}
			if err != nil {
				log.Fatalf("Error while receiving from biDi client: %v", err)
			}

			fmt.Printf("Received: %v\n", res.GetResult())
		}
	}()

	// block until everything is done
	<-doneChannel
}

func doUnaryCallWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting UnaryWithDealine RPC...", timeout)

	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Angel",
			LastName:  "Dionisio",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// make call to GreetWithDealine, if we get an error from this call
	// we take a look at the error, if the error is a 'DeadlineExceeded' error
	// we print a custom message
	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {

		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Dealine was exceeded")
			} else {
				fmt.Printf("Unexepected error: %v", statusErr)
			}

		} else {
			log.Fatalf("error while calling GreetWithDealine RPC: %v", err)
		}
		// return in case of any err
		return
	}

	log.Printf("response from server GreetWithDealine: %v\n", res.Result)
}
