package main

import (
	"context"
	"tolling/types"
)

type Aggregator interface {
	Aggregate(context.Context, types.Distance) error
}

type Storer interface {
	Store(context.Context, types.Distance) error
}

type InvoiceAggregator struct {
	store Storer
}

func (ia *InvoiceAggregator) Aggregate(ctx context.Context, d types.Distance) error {
	return ia.store.Store(ctx, d)
}

func NewInvoceAggregator(store Storer) *InvoiceAggregator {
	return &InvoiceAggregator{store}
}

// I thing it's going to be Invoicer specific
