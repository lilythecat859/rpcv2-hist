package grpc

import (
	"context"
	"crypto/tls"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/lilythecat859/rpcv2-hist/internal/api/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.HistoricalClient
}

func NewClient(addr string, insecureOpt bool) (*Client, error) {
	var opts []grpc.DialOption
	if insecureOpt {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		tc := credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})
		opts = append(opts, grpc.WithTransportCredentials(tc))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:   conn,
		client: pb.NewHistoricalClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetBlock(ctx context.Context, slot uint64, commitment string) ([]byte, error) {
	resp, err := c.client.GetBlock(ctx, &GetBlockRequest{Slot: slot, Commitment: commitment})
	if err != nil {
		return nil, err
	}
	return resp.Raw, nil
}

func (c *Client) GetTransaction(ctx context.Context, sig string, commitment string) ([]byte, error) {
	resp, err := c.client.GetTransaction(ctx, &GetTransactionRequest{Signature: sig, Commitment: commitment})
	if err != nil {
		return nil, err
	}
	return resp.Raw, nil
}