package aggregator

import (
	"context"
	"tolling/types"
)

// factory pattern
type Middleware func(Service) Service

// logging middleware
type loggingMiddleware struct {
	next Service
}

func newLoggingMiddleware() Middleware {
	return func(next Service) Service {
		return loggingMiddleware{next: next}
	}
}

func (mw loggingMiddleware) Aggregate(ctx context.Context, d types.Distance) error {
	return mw.next.Aggregate(ctx, d)
}

func (mw loggingMiddleware) GetInvoice(ctx context.Context, OBUDID int) (types.Invoice, error) {
	return mw.next.GetInvoice(ctx, OBUDID)
}

func (mw loggingMiddleware) Shutdown() error {
	return nil
}

// instrumentation middleware
type instrumentationMiddleware struct {
	next Service
}

func newInstrumentationMiddleware() Middleware {
	return func(next Service) Service {
		return instrumentationMiddleware{next: next}
	}
}

func (mw instrumentationMiddleware) Aggregate(ctx context.Context, d types.Distance) error {
	return mw.next.Aggregate(ctx, d)
}

func (mw instrumentationMiddleware) GetInvoice(ctx context.Context, OBUDID int) (types.Invoice, error) {
	return mw.next.GetInvoice(ctx, OBUDID)
}

func (mw instrumentationMiddleware) Shutdown() error {
	return nil
}
