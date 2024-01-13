package main

import (
	"context"
	"tolling/common"
)

type LogMiddleware[T any] struct {
	next   Consumer[T]
	logger common.Logger
}

func NewLogMiddleware[T any](next Consumer[T], logger common.Logger) *LogMiddleware[T] {
	return &LogMiddleware[T]{next, logger}
}

func (lw *LogMiddleware[T]) Consume(ctx context.Context) chan EventPayload[T] {
	defer lw.logger.New().Info("started to consume from kafka")
	return lw.next.Consume(ctx)
}
