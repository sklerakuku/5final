// internal/grpc/client.go
package grpc

import (
	"context"
	"time"

	"github.com/sklerakuku/5final/internal/config"
	pb "github.com/sklerakuku/5final/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.CalculatorClient
	config *config.Config
	sem    chan struct{}
}

func NewClient(addr string, cfg *config.Config) (*Client, error) {
	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:   conn,
		client: pb.NewCalculatorClient(conn),
		config: cfg,
		sem:    make(chan struct{}, cfg.ComputingPower),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Calculate(arg1, arg2 float64, op string) (float64, error) {
	c.sem <- struct{}{}
	defer func() { <-c.sem }()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := c.client.Calculate(ctx, &pb.Task{
		Arg1:      arg1,
		Arg2:      arg2,
		Operation: op,
	})
	if err != nil {
		return 0, err
	}

	if res.Error != nil {
		return 0, res.Error
	}

	return res.Value, nil
}
