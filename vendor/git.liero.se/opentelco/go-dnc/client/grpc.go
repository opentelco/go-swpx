package client

import (
	"context"
	"time"

	"git.liero.se/opentelco/go-dnc/models/protobuf/dispatcher"
	"git.liero.se/opentelco/go-dnc/models/protobuf/transport"
	"google.golang.org/grpc"
)

const (
	address = "127.0.0.1:50051"
)

//GRPCClient Calls the client.
type GRPCClient struct {
	conn *grpc.ClientConn
	// client interface
	client dispatcher.DispatcherClient
}

// New creates a new Client
func New(addr string) (Client,error){
	var err error
	client := &GRPCClient{}
	// Set up a connection to the server.
	client.conn, err = grpc.Dial(addr, grpc.WithInsecure(),grpc.WithTimeout(time.Duration(time.Second*5)))
	if err != nil {
		return nil,err
	}

	client.client = dispatcher.NewDispatcherClient(client.conn)
	return client, nil

}

func (s *GRPCClient) Close() error {
	return s.conn.Close()
}

// Ping sends a ping request..
func (c *GRPCClient) Ping() (*dispatcher.PingReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	r, err := c.client.Ping(ctx, &dispatcher.PingRequest{CallerId: "1234"})

	return r, err
}

func (c *GRPCClient) Put(req *transport.Message) (*transport.Message, error) {
	return c.client.Put(context.Background(), req)
}
