package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	pb "github.com/andreasatle/Udemy/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *pb.SumRequest) (*pb.SumResponse, error) {
	fmt.Printf("Sum invoked on server: %v\n", req)
	firstNumber := req.GetFirstNumber()
	secondNumber := req.GetSecondNumber()
	result := firstNumber + secondNumber
	res := &pb.SumResponse{
		SumResult: result,
	}
	return res, nil
}

func (*server) Prime(req *pb.PrimeRequest, stream pb.CalculatorService_PrimeServer) error {
	fmt.Printf("Prime invoked on server: %v\n", req)
	for divisor, number := int64(2), req.GetNumber(); number > 1 && divisor <= number; {
		if number%divisor == 0 {
			number /= divisor
			res := &pb.PrimeResponse{
				PrimeResult: divisor,
			}
			stream.Send(res)
			//time.Sleep(1000 * time.Millisecond)
		} else {
			divisor++
		}
	}
	return nil
}

func (*server) Average(stream pb.CalculatorService_AverageServer) error {
	fmt.Println("Average function invoked on server")
	sum := int64(0)
	nTerms := int64(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float64(sum) / float64(nTerms)
			return stream.SendAndClose(&pb.AverageResponse{
				Average: average,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		fmt.Printf("Received number: %v\n", req.GetNumber())
		sum += req.GetNumber()
		nTerms++
	}
}

func (*server) Max(stream pb.CalculatorService_MaxServer) error {
	fmt.Println("Max function invoked on server")
	maxNum := int64(math.MinInt64)

	for {
		req, recvErr := stream.Recv()
		if recvErr == io.EOF {
			return nil
		}
		if recvErr != nil {
			log.Fatalf("Error while reading client data: %v\n", recvErr)
			return recvErr
		}

		// Check if num is maximum so far
		num := req.GetNumber()
		fmt.Printf("Max received number: %v\n", num)
		if num <= maxNum {
			continue
		}
		maxNum = num

		sendErr := stream.Send(&pb.MaxResponse{
			Max: maxNum,
		})
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v\n", sendErr)
			return sendErr
		}

	}
}

func (*server) SquareRoot(ctx context.Context, req *pb.SquareRootRequest) (*pb.SquareRootResponse, error) {
	fmt.Printf("SquareRoot invoked on server: %v\n", req)
	number := req.GetNumber()
	fmt.Printf("SquareRoot received number: %v\n", number)
	if number < 0.0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative argument: %v", number),
		)
	}
	return &pb.SquareRootResponse{
		SquareRoot: math.Sqrt(number),
	}, nil
}

func main() {
	fmt.Println("Hello world, from Calculator server!")

	// Listen to port (some default grpc)
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new server
	s := grpc.NewServer()

	// Register the server at pb
	pb.RegisterCalculatorServiceServer(s, &server{})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
