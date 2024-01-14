package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tolling/aggregator/client"
	"tolling/common"
	"tolling/types"
)

const (
	kafkaTopic      = "test-topic"
	kafkaAddrs      = "127.0.0.1:29092"
	aggHttpEndpoint = "http://127.0.0.1:3000"
	aggGrpcEndpoint = "127.0.0.1:50051"
)

func main() {
	logger := common.NewCustomLogger()
	l := logger.New()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-shutdown
		l.Info("received signal to cancel")
		cancel()
	}()

	con, err := NewKafkaConsumer[types.OBUData](KafkaConsumerConfig{
		addrs:   kafkaAddrs,
		topics:  []string{kafkaTopic},
		groupID: "groupID",
	},
		logger,
	)

	if err != nil {
		log.Fatal(err)
	}

	aggHTTP := client.NewAggregatorHttpClient(aggHttpEndpoint, logger)
	aggGRPC, err := client.NewAggregatorGrpcClient(aggGrpcEndpoint, logger)
	if err != nil {
		log.Fatalln("GRPS is not connected")
	}

	calc := NewDistanceCalculator(con, aggGRPC, logger)

	calc.Run(ctx)
}
