package main

import (
	"context"
	"tolling/common"
	"tolling/types"
)

type LogMiddleware struct {
	next   Producer
	logger common.Logger
}

func NewLogMiddleware(next Producer, logger common.Logger) *LogMiddleware {
	return &LogMiddleware{next, logger}
}

func (lw *LogMiddleware) Produce(ctx context.Context, data *types.OBUData) (err error) {
	defer func() {
		l := lw.logger.New()

		if traceID, ok := ctx.Value(types.KeyTraceID).(string); ok {
			l.WithTraceID(traceID)
		}

		l.WithOBUID(data.OBUID)

		if err != nil {
			l.WithError(err).Error("produce failed")
		} else {
			l.Info("produce succeeded")
		}

	}()

	err = lw.next.Produce(ctx, data)
	return err
}
