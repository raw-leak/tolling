package main

import (
	"context"
	"net"
	"tolling/common"
	"tolling/types"

	"google.golang.org/grpc"
)

type GrpcTransport struct {
	types.UnimplementedAggregatorServer

	aggregator Aggregator
	logger     common.Logger
	server     *grpc.Server
}

func NewGrpcTransport(aggregator Aggregator, logger common.Logger) *GrpcTransport {
	return &GrpcTransport{aggregator: aggregator, logger: logger}
}

func (t *GrpcTransport) Aggregate(ctx context.Context, req *types.AggregateReq) (*types.None, error) {
	l := t.logger.New()

	distance := types.Distance{
		Value: req.Value,
		OBUID: int(req.ObuId),
		Unix:  req.Unix,
	}

	if err := t.aggregator.Aggregate(ctx, distance); err != nil {
		l.Error("aggregate failed")
		return &types.None{}, err
	}

	return &types.None{}, nil
}

func (t *GrpcTransport) StartGrpcServer(listenAddr string) error {
	l := t.logger.New()

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		l.WithError(err).Error("failed to listen GRPC")
		return err
	}

	t.server = grpc.NewServer()
	types.RegisterAggregatorServer(t.server, t)

	l.Infof("starting gRPC server on %s", listenAddr)
	return t.server.Serve(ln)
}

func (t *GrpcTransport) Shutdown() {
	t.logger.New().Infof("stopping gRPC server")

	if t.server != nil {
		t.server.GracefulStop()
	}
}
