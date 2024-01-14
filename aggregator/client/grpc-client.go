package client

import (
	"context"
	"time"
	"tolling/common"
	"tolling/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

const (
	timeout             = 5 * time.Second
	keepAliveTime       = 30 * time.Second
	keepAliveTimeout    = 10 * time.Second
	keepAlivePermit     = false
	retryMaxAttempts    = 3
	retryInitialBackoff = 1 * time.Second
	retryMaxBackoff     = 5 * time.Second
)

type AggregatorGrpcClient struct {
	grpcAddr string
	logger   common.Logger
	conn     *grpc.ClientConn
	types.AggregatorClient
}

func NewAggregatorGrpcClient(grpcAddr string, logger common.Logger) (*AggregatorGrpcClient, error) {
	opts := []grpc.DialOption{
		// grpc.WithBlock(), // block until the connection is established
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                keepAliveTime,
			Timeout:             keepAliveTimeout,
			PermitWithoutStream: keepAlivePermit,
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(grpcAddr, opts...)
	if err != nil {
		return nil, err
	}

	ac := types.NewAggregatorClient(conn)

	logger.New().Infof("gRPC AggregatorClient has been successfully connected to %s", grpcAddr)
	return &AggregatorGrpcClient{grpcAddr: grpcAddr, logger: logger, conn: conn, AggregatorClient: ac}, nil
}

func (c *AggregatorGrpcClient) Aggregate(ctx context.Context, d types.Distance) error {
	req := types.AggregateReq{Value: d.Value, Unix: d.Unix, OBUID: int32(d.OBUID)}
	_, err := c.AggregatorClient.Aggregate(ctx, &req)
	if err != nil {
		return err
	}

	return nil
}

func (c *AggregatorGrpcClient) GetInvoice(ctx context.Context, OBUID int) (types.Invoice, error) {
	req := types.GetInvoiceReq{OBUID: int32(OBUID)}
	invRes, err := c.AggregatorClient.GetInvoice(ctx, &req)
	if err != nil {
		return types.Invoice{}, err
	}

	return types.Invoice{
		OBUID:         int(invRes.OBUID),
		TotalDistance: float64(invRes.TotalDistance),
		TotalAmount:   float64(invRes.TotalAmount),
	}, nil
}

func (c *AggregatorGrpcClient) Shutdown() error {
	c.logger.New().Infof("AggregatorClient is disconnecting")

	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
