package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"github.com/tkircsi/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct{}

func main() {
	log.Println("Calculator Server started...")
	l, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	log.Fatal(s.Serve(l))
}

func (s *server) Sum(ctx context.Context, r *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Println("Received Sum RPC")
	time.Sleep(500 * time.Millisecond)
	return &calculatorpb.SumResponse{
		Result: r.GetN1() + r.GetN2(),
	}, nil
}

func (s *server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	log.Println("Received Average RPC")
	var sum, count int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.AverageResponse{
				Average: float32(sum) / float32(count),
			})
		}
		if err != nil {
			log.Println(err)
			return err
		}
		n := req.GetN1()
		sum += n
		count++
	}
}

func (s *server) GetPrimes(r *calculatorpb.PrimeDecompRequest, stream calculatorpb.CalculatorService_GetPrimesServer) error {
	log.Println("Received GetPrimes RPC")
	num := r.GetNumber()
	var k int64 = 2
	for num > 1 {
		if num%k == 0 {
			time.Sleep(500 * time.Millisecond)
			stream.Send(&calculatorpb.PrimeDecompResponse{
				PrimeNumber: k,
			})
			num = num / k
		} else {
			k++
		}
	}
	return nil
}

func (s *server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	log.Println("Received FindMaximum RPC")
	var max int32 = 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("error while receiving request from client: %v\n", err)
			return err
		}
		n := req.GetN()
		log.Printf("Received from cient: %d\n", n)
		if max < n || max == 0 {
			max = n
			err = stream.Send(&calculatorpb.MaximumResponse{
				Maximum: max,
			})
			if err != nil {
				log.Printf("error while sending response to client: %v\n", err)
				return err
			}
		}
	}
}

func (s *server) Sqrt(ctx context.Context, r *calculatorpb.SqrtRequest) (*calculatorpb.SqrtResponse, error) {
	log.Println("Received Sqrt RPC")
	n := r.GetN()
	if n < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Received a negative number: %d", n))
	}
	return &calculatorpb.SqrtResponse{Root: math.Sqrt(float64(n))}, nil
}
