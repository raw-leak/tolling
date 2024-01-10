package main

import (
	"context"
	"encoding/json"
	"tolling/common"
	"tolling/types"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaProducerConfig struct {
	addrs string
	topic string
}

type KafkaProducer struct {
	kp     *kafka.Producer
	cfg    KafkaProducerConfig
	logger common.Logger
}

func NewKafkaProducer(cfg KafkaProducerConfig, logger common.Logger) (Producer, error) {
	l := logger.New()

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": cfg.addrs})
	if err != nil {
		l.Error("kafka-producer failed to start")
		return nil, err
	}

	_, err = p.GetMetadata(nil, false, 5000)
	if err != nil {
		l.Errorf("kafka-producer failed to connect: \n%v", err)
	} else {
		l.Info("kafka-producer successfully connected")
	}

	// delivery report handler for produced messages
	// go func() {
	// 	for e := range p.Events() {
	// 		switch ev := e.(type) {
	// 		case *kafka.Message:
	// 			if ev.TopicPartition.Error != nil {
	// 				fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
	// 			} else {
	// 				fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
	// 			}
	// 		}
	// 	}
	// }()

	return &KafkaProducer{
		p,
		cfg,
		logger,
	}, nil
}

func (p *KafkaProducer) Produce(ctx context.Context, data *types.OBUData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	headers := []kafka.Header{}
	if traceID, ok := ctx.Value(types.KeyTraceID).(string); ok {
		headers = append(headers, kafka.Header{
			Key:   string(types.KeyTraceID),
			Value: []byte(traceID),
		})
	}

	return p.kp.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.cfg.topic,
			Partition: kafka.PartitionAny,
		},
		Value:   b,
		Headers: headers,
	}, nil)
}
