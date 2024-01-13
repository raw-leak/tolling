package main

import (
	"context"
	"encoding/json"
	"fmt"
	"tolling/common"
	"tolling/types"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Logger interface {
	Info(msg string, fields map[string]any)
	Error(msg string, fields map[string]any)
}

type EventPayload[T any] struct {
	data T
	ctx  context.Context
}

type KafkaConsumerConfig struct {
	groupID string
	addrs   string
	topics  []string
}

type KafkaConsumer[T any] struct {
	con    *kafka.Consumer
	cfg    KafkaConsumerConfig
	logger common.Logger
}

func NewKafkaConsumer[T any](cfg KafkaConsumerConfig, logger common.Logger) (Consumer[T], error) {
	con, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.addrs,
		"group.id":          cfg.groupID,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	if err := con.SubscribeTopics(cfg.topics, nil); err != nil {
		return nil, err
	}

	return &KafkaConsumer[T]{con, cfg, logger}, nil
}

func (kc *KafkaConsumer[T]) Consume(ctx context.Context) chan EventPayload[T] {
	conCh := make(chan EventPayload[T])

	go func() {
		defer kc.con.Close()
		defer close(conCh)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("stopping the loop")
				return
			default:
				var data T
				msg, err := kc.con.ReadMessage(-1)
				l := kc.logger.New()

				if err != nil {
					if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.IsTimeout() {
						continue
					}
					l.Errorf("read kafka message error:  %v\n", err)
					continue
				}

				ctx := context.Background()
				if len(msg.Headers) > 0 {
					traceID := string(msg.Headers[0].Value)
					ctx = context.WithValue(ctx, types.KeyTraceID, traceID)
					l.WithTraceID(traceID)
				}

				if err := json.Unmarshal(msg.Value, &data); err != nil {
					l.Errorf("failed to marshal consumed OBU data from kafka with error: \n%v", err)
					continue
				}

				l.Info("consumed OBU data")

				v, o := ctx.Value(types.KeyTraceID).(string)
				fmt.Println("before sendit the evnet to channel ->>> ", v, o)

				conCh <- EventPayload[T]{data, ctx}
			}
		}
	}()

	return conCh
}
