// cmd/agent/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/sklerakuku/5final/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCalculatorServer
}

func (s *server) Calculate(ctx context.Context, req *pb.Task) (*pb.Result, error) {
	var result float64
	var errMsg error

	switch req.Operation {
	case "+":
		result = req.Arg1 + req.Arg2
	case "-":
		result = req.Arg1 - req.Arg2
	case "*":
		result = req.Arg1 * req.Arg2
	case "/":
		if req.Arg2 == 0 {
			errMsg = fmt.Errorf("division by zero")
		} else {
			result = req.Arg1 / req.Arg2
		}
	default:
		errMsg = fmt.Errorf("unsupported operation")
	}

	return &pb.Result{
		Value: result,
		Error: errMsg,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCalculatorServer(s, &server{})

	log.Println("Agent server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
