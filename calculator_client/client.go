package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/tkircsi/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	// Unary RPC
	SumTest(c)
	fmt.Println()
	fmt.Println()
	time.Sleep(2 * time.Second)

	// Client streaming
	AverageTest(c)
	fmt.Println()
	fmt.Println()
	time.Sleep(2 * time.Second)

	// Server streaming
	PrimeDecompTest(c)
	fmt.Println()
	fmt.Println()
	time.Sleep(2 * time.Second)

	// BiDi RPC
	MaximumTest(c)
	fmt.Println()
	fmt.Println()
	time.Sleep(2 * time.Second)

	// Unary RPC
	SqrtTest(c) // Test error and error handling
	fmt.Println()
	fmt.Println()
	time.Sleep(2 * time.Second)
}

func SqrtTest(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Sqrt Unary RPC...")
	fmt.Println(strings.Repeat("-", 60))
	n := int32(36)

	reqNoError := &calculatorpb.SqrtRequest{N: n}
	fmt.Printf("Sending SqrtRequest: %s\n", reqNoError.String())
	resp, err := c.Sqrt(context.Background(), reqNoError)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received SqrtResponse: %s\n", resp)

	n = int32(-4)
	reqNoError = &calculatorpb.SqrtRequest{N: n}
	fmt.Printf("Sending SqrtRequest: %s\n", reqNoError.String())
	resp, err = c.Sqrt(context.Background(), reqNoError)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received SqrtResponse: %s\n", resp)
}

func MaximumTest(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a MaximumTest BiDi RPC...")
	fmt.Println(strings.Repeat("-", 60))
	reqs := []int32{-2, -3, 13, 6, 8, 17, -23, -2, 33}
	ch := make(chan struct{})
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Nums: %v\n", reqs)

	go func() {

		for _, r := range reqs {
			req := &calculatorpb.MaximumRequest{N: r}
			err := stream.Send(req)
			if err != nil {
				log.Println(err)
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		fmt.Println("Receiving:")
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println(err)
				break
			}
			fmt.Print(resp.GetMaximum(), ", ")
		}
		close(ch)
	}()

	<-ch
	fmt.Println()
}

func SumTest(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a SumTest Unary RPC...")
	fmt.Println(strings.Repeat("-", 60))
	req := calculatorpb.SumRequest{
		N1: 10,
		N2: 100,
	}
	fmt.Printf("Sending SumRequest: %s\n", req.String())
	res, err := c.Sum(context.Background(), &req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received SumResponse: %s\n", res.String())
}

func AverageTest(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a AverageTest Client Streaming RPC...")
	fmt.Println(strings.Repeat("-", 60))
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	nums := []int32{2, 4, 6, 7}
	fmt.Printf("Numbers: %v\n", nums)

	for _, i := range nums {
		req := &calculatorpb.AverageRequest{N1: int32(i)}
		fmt.Printf("Sending AverageRequest: %s\n", req.String())
		time.Sleep(500 * time.Millisecond)
		err = stream.Send(req)
		if err != nil {
			log.Fatal(err)
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received AverageResponse: %s\n", resp.String())
}

func PrimeDecompTest(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a PrimeDecompTest Server Streaming RPC...")
	fmt.Println(strings.Repeat("-", 60))
	req := &calculatorpb.PrimeDecompRequest{
		// Number: 566732120123,
		Number: 233456,
		//Number: 5,
	}
	fmt.Printf("Sending PrimeDecompRequest: %s\n", req.String())
	stream, err := c.GetPrimes(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Receiving:")
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		i := resp.GetPrimeNumber()
		fmt.Print(i, ",")
	}
	fmt.Println()
}
