package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/andreasatle/Udemy/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello world, from Calculator client!")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Client could not connect to server: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	for i := 0; i < 10; i++ {
		doUnary(c)
		doServerStreaming(c)
		doClientStreaming(c, []int64{7, 5, 2, 1})
		doClientStreaming(c, []int64{1, 3, 5, 7, 9})
		doBiDiStreaming(c, []int64{3, 6, 4, 1, 9, 4, 6, 2, 11, 47, 3})
		doBiDiStreaming(c, []int64{4, 1, 9, 4, 6, 2, 11, 47, 3, 53, 1, 2, 3, 4, 5})
		doErrorUnary(c, 10.0)
		doErrorUnary(c, -2.0)
	}
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting a Unary RPC...")

	firstNumber := int64(10)
	secondNumber := int64(3)

	req := &calculatorpb.SumRequest{
		FirstNumber:  firstNumber,
		SecondNumber: secondNumber,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error Sum RPC: %v", err)
	}

	log.Printf("Response from Sum: %d+%d=%v", firstNumber, secondNumber, res.SumResult)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting a Server Streaming RPC...")

	number := int64(12345678)

	req := &calculatorpb.PrimeRequest{
		Number: number,
	}

	stream, err := c.Prime(context.Background(), req)
	if err != nil {
		log.Fatalf("Error Prime RPC: %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// End of the stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from Prime: %v", msg.GetPrimeResult())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient, numbers []int64) {
	fmt.Println("Starting a Client Streaming RPC...")

	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("Error while calling Average: %v", err)
	}

	for _, number := range numbers {
		req := &calculatorpb.AverageRequest{
			Number: int64(number),
		}
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from Average: %v", err)
	}
	fmt.Printf("Response from Average: %v\n", res)
}

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient, numbers []int64) {
	fmt.Println("Starting a Bi-Di Streaming RPC...")

	var maxNum int64

	stream, err := c.Max(context.Background())
	if err != nil {
		log.Fatalf("Error calling Max RPC: %v\n", err)
	}

	waitc := make(chan struct{})
	go func() {
		for _, num := range numbers {
			req := &calculatorpb.MaxRequest{
				Number: num,
			}
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving data: %v\n", err)
				break
			}
			maxNum = res.GetMax()
			fmt.Printf("Received: %v\n", maxNum)
		}
		close(waitc)
	}()

	<-waitc
	fmt.Printf("Response from Max: %v\n", maxNum)
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient, number float64) {
	fmt.Println("Starting a Unary RPC with Error Handling...")
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: number})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
			}
		} else {
			log.Fatalf("Big Error calling SquareRoot: %v\n", err)
		}
	}
	fmt.Printf("Response from SquareRoot: %v\n", res.GetSquareRoot())
}
