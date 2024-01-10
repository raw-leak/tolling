package main

import (
	"context"

	"github.com/sirupsen/logrus"
)

type LogMiddleware[T any] struct {
	next Consumer[T]
}

func NewLogMiddleware[T any](next Consumer[T]) *LogMiddleware[T] {
	return &LogMiddleware[T]{next}
}

func (lw *LogMiddleware[T]) Consume(ctx context.Context) chan T {
	logrus.Info("started consuming from topics: []")
	return lw.next.Consume(ctx)
}
