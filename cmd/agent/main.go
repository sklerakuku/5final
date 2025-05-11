package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/sklerakuku/5final/internal/config"
	pb "github.com/sklerakuku/5final/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCalculatorServer
	config *config.Config
}

func (s *server) Calculate(ctx context.Context, req *pb.Task) (*pb.Result, error) {
	var result float64
	var errMsg error

	switch req.Operation {
	case "+":
		time.Sleep(time.Duration(s.config.TimeAdditionMS) * time.Millisecond)
		result = req.Arg1 + req.Arg2
	case "-":
		time.Sleep(time.Duration(s.config.TimeSubtractionMS) * time.Millisecond)
		result = req.Arg1 - req.Arg2
	case "*":
		time.Sleep(time.Duration(s.config.TimeMultiplicationMS) * time.Millisecond)
		result = req.Arg1 * req.Arg2
	case "/":
		time.Sleep(time.Duration(s.config.TimeDivisionMS) * time.Millisecond)
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
	cfg := config.Load()

	opts := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(uint32(cfg.ComputingPower)),
	}

	srv := grpc.NewServer(opts...)
	pb.RegisterCalculatorServer(srv, &server{config: cfg})

	lis, err := net.Listen("tcp", cfg.GRPCAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Agent server listening on %s with computing power %d", cfg.GRPCAddress, cfg.ComputingPower)
	log.Printf("Operation timings - Add: %dms, Sub: %dms, Mul: %dms, Div: %dms",
		cfg.TimeAdditionMS, cfg.TimeSubtractionMS, cfg.TimeMultiplicationMS, cfg.TimeDivisionMS)

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
