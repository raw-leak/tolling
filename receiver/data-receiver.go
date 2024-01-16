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
	prod   Producer
	logger common.Logger
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin function should be carefully configured to avoid cross-site request forgery (CSRF) attacks.
}

func NewDataReceiver(prod Producer, logger common.Logger) (*DataReceiver, error) {
	return &DataReceiver{
		msgCh:  make(chan types.OBUData, 128),
		prod:   prod,
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

	go dr.wsReceiveLoop(conn)
}

func (dr *DataReceiver) wsReceiveLoop(conn *websocket.Conn) {
	wsl := dr.logger.New()
	wsl.Info("new OBU connected")
	for {
		l := dr.logger.New()

		var data types.OBUData
		if err := conn.ReadJSON(&data); err != nil {
			l.Errorf("failed to read from ws: \n%v", err)

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				l.Errorf("ws connection closed unexpectedly: \n%v", err)
			} else {
				l.Infof("ws connection closed: \n%v", err)
				conn.WriteMessage(websocket.CloseMessage, []byte{})

			}
			break
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

	wsl.Info("new OBU disconnected")

}
