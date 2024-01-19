package main

import (
	"context"
	"tolling/common"
	"tolling/types"
)

type LogMiddleware struct {
	next   Aggregator
	logger common.Logger
}

func NewLogMiddleware(next Aggregator, logger common.Logger) Aggregator {
	return &LogMiddleware{next, logger}
}

func (lw *LogMiddleware) Aggregate(ctx context.Context, d types.Distance) (err error) {
	defer func() {
		l := lw.logger.New()
		if traceID, ok := ctx.Value(types.KeyTraceID).(string); ok {
			l.WithTraceID(traceID)
		}

		if err != nil {
			l.WithError(err).Error("aggregate failed")
		} else {
			l.WithOBUID(d.OBUID).Info("aggregate succeeded")
		}
	}()

	err = lw.next.Aggregate(ctx, d)
	return
}

func (lw *LogMiddleware) GetInvoice(ctx context.Context, OBUID int) (invoice types.Invoice, err error) {
	defer func() {
		l := lw.logger.New()
		if traceID, ok := ctx.Value(types.KeyTraceID).(string); ok {
			l.WithTraceID(traceID)
		}

		if err != nil {
			l.WithError(err).Error("invoice calculation failed")
		} else {
			l.WithOBUID(OBUID).Info("invoice calculation succeeded")
		}
	}()

	invoice, err = lw.next.GetInvoice(ctx, OBUID)
	return
}

func (lw *LogMiddleware) Shutdown() error {
	return nil
}
