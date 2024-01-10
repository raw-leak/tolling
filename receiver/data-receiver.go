package main

import (
	"context"
	"net/http"
	"tolling/common"
	"tolling/types"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Producer interface {
	Produce(context.Context, *types.OBUData) error
}

type DataReceiver struct {
	msgCh  chan types.OBUData
	conn   *websocket.Conn
	prod   Producer
	logger common.Logger
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin function should be carefully configured to avoid cross-site request forgery (CSRF) attacks.
}

func NewDataReceiver(logger common.Logger) (*DataReceiver, error) {
	p, err := NewKafkaProducer(KafkaProducerConfig{
		addrs: kafkaAddrs,
		topic: kafkaTopic,
	}, logger)
	if err != nil {
		return nil, err
	}
	p = NewLogMiddleware(p, logger)

	return &DataReceiver{
		msgCh:  make(chan types.OBUData, 128),
		prod:   p,
		logger: logger,
	}, nil
}

func (dr *DataReceiver) handler(w http.ResponseWriter, r *http.Request) {
	l := dr.logger.New()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		l.Error("http upgrade to ws failed")
		return
	}

	dr.conn = conn

	go dr.wsReceiveLoop()
}

func (dr *DataReceiver) wsReceiveLoop() {
	dr.logger.New().Info("new OBU connected")

	for {
		l := dr.logger.New()

		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			l.Errorf("failed to read from ws: \n%v", err)
			continue
		}

		traceID := uuid.New().String()
		ctx := context.WithValue(context.Background(), types.KeyTraceID, traceID)

		l.WithTraceID(traceID)

		if err := dr.prod.Produce(ctx, &data); err != nil {
			l.Errorf("send to produce failed: \n%v", err)
		} else {
			l.Info("send to produce succeeded")
		}

	}
}
