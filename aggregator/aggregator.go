package main

import (
	"context"
	"errors"
	"tolling/types"
)

type Aggregator interface {
	Aggregate(context.Context, types.Distance) error
	GetInvoice(context.Context, int) (types.Invoice, error)
}

type Storer interface {
	Store(context.Context, types.Distance) error
	Read(context.Context, int) (float64, bool)
}

type InvoiceAggregator struct {
	store Storer
}

func (ia InvoiceAggregator) Aggregate(ctx context.Context, d types.Distance) error {
	return ia.store.Store(ctx, d)
}

func (ia InvoiceAggregator) GetInvoice(ctx context.Context, OBUDID int) (types.Invoice, error) {
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

func (ia InvoiceAggregator) calculateInvoice(ctx context.Context, totalDistance float64) float64 {
	return totalDistance * 1.5
}

func NewInvoiceAggregator(store Storer) *InvoiceAggregator {
	return &InvoiceAggregator{store}
}

// I thing it's going to be Invoicer specific
