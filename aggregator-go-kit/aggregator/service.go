package aggregator

import (
	"context"
	"errors"
	"tolling/types"
)

type Service interface {
	Aggregate(context.Context, types.Distance) error
	GetInvoice(context.Context, int) (types.Invoice, error)
	Shutdown() error
}

type storer interface {
	Store(context.Context, types.Distance) error
	Read(context.Context, int) (float64, bool)
}

type invoiceAggregator struct {
	store storer
}

func (ia invoiceAggregator) Aggregate(ctx context.Context, d types.Distance) error {
	return ia.store.Store(ctx, d)
}

func (ia invoiceAggregator) GetInvoice(ctx context.Context, OBUDID int) (types.Invoice, error) {
	totalDistance, ok := ia.store.Read(ctx, OBUDID)
	if !ok {
		return types.Invoice{}, errors.New("OBU data not found")
	}

	totalAmount := ia.calculateInvoice(ctx, totalDistance)
	return types.Invoice{
		TotalDistance: totalAmount,
		TotalAmount:   totalAmount,
		OBUID:         OBUDID,
	}, nil
}

func (ia invoiceAggregator) Shutdown() error {
	return nil
}

func (ia invoiceAggregator) calculateInvoice(ctx context.Context, totalDistance float64) float64 {
	return totalDistance * 1.5
}

func newInvoiceAggregator(store storer) Service {
	return &invoiceAggregator{store: store}
}

// NewInvoiceAggregator will construct a complete microservice
// with Logging and Instrumentation middlewares
func NewInvoiceAggregatorService(store storer) Service {
	var ias Service

	{
		ias = newInvoiceAggregator(newMemoryStorer())
		ias = newLoggingMiddleware()(ias)
		ias = newInstrumentationMiddleware()(ias)
	}

	return ias
}
