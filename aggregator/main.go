package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"tolling/common"
)

var httpAddr = flag.String("httpAddr", ":3000", "the listen address of the HTTP server")
var grpcAddr = flag.String("grpcAddr", ":50051", "the listen address of the gRPC server")

func main() {
	flag.Parse()
	logger := common.NewCustomLogger()
	l := logger.New()

	var aggregator Aggregator

	memoryStore := NewMemoryStorer()

	aggregator = NewInvoiceAggregator(memoryStore)
	aggregator = NewLogMiddleware(aggregator, logger)

	httpTransport := NewHttpTransport(aggregator, logger)
	grpcTransport := NewGrpcTransport(aggregator, logger)

	errChan := make(chan error, 2)

	go func() {
		if err := httpTransport.StartHttpServer(*httpAddr); err != nil {
			errChan <- err
		}
	}()

	go func() {
		if err := grpcTransport.StartGrpcServer(*grpcAddr); err != nil {
			errChan <- err
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		l.Info("Shutdown signal received")
	case err := <-errChan:
		l.Errorf("Server error: %v", err)
	}

	if err := httpTransport.Shutdown(context.Background()); err != nil {
		l.Errorf("HTTP server shutdown error: %v", err)
	} else {
		l.Info("HTTP server gracefully shutdown")
	}

	grpcTransport.Shutdown()
	l.Info("gRPC server gracefully shutdown")
}
