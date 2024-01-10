package main

import (
	"encoding/json"
	"net/http"
	"tolling/common"
	"tolling/types"

	"context"
)

type HttpTransport struct {
	aggregator Aggregator
	logger     common.Logger
}

func NewHttpTransport(aggregator Aggregator, logger common.Logger) *HttpTransport {
	return &HttpTransport{aggregator, logger}
}

func (t *HttpTransport) handleAggregate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	l := t.logger.New()

	ctx := context.Background()
	if traceIds, ok := r.Header[types.TraceIDHeader]; ok && len(traceIds) > 0 {
		ctx = context.WithValue(ctx, types.KeyTraceID, traceIds[0])
		l.WithTraceID(traceIds[0])
	}

	var distance types.Distance
	if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
		l.Error("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := t.aggregator.Aggregate(ctx, distance); err != nil {
		l.Error("aggregate failed")
		http.Error(w, "aggregate failed", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	l.Info("aggregate processed")
}
