package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tolling/common"
	"tolling/types"

	"context"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HttpTransport struct {
	aggregator Aggregator
	logger     common.Logger
	server     *http.Server
}

func NewHttpTransport(aggregator Aggregator, logger common.Logger) *HttpTransport {
	return &HttpTransport{aggregator, logger, nil}
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
}

func (t *HttpTransport) handleGetInvoice(w http.ResponseWriter, r *http.Request) {
	l := t.logger.New()
	ctx := context.Background()

	if traceIds, ok := r.Header[types.TraceIDHeader]; ok && len(traceIds) > 0 {
		ctx = context.WithValue(ctx, types.KeyTraceID, traceIds[0])
		l.WithTraceID(traceIds[0])
	}

	strOBUID := r.URL.Query().Get("obuid")
	if len(strOBUID) < 1 {
		http.Error(w, "no obuid provided", http.StatusBadRequest)
		return
	}

	OBUDID, err := strconv.Atoi(strOBUID)
	if err != nil {
		http.Error(w, "Invalid obuid query", http.StatusBadRequest)
		return
	}

	distance, err := t.aggregator.GetInvoice(ctx, OBUDID)
	if err != nil {
		http.Error(w, "invoice calculation failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"distance": distance})
}

func (t *HttpTransport) StartHttpServer(httpAddr string) error {
	t.server = &http.Server{Addr: httpAddr, Handler: nil}

	http.HandleFunc("/aggregate", t.handleAggregate)
	http.HandleFunc("/invoice", t.handleGetInvoice)

	http.Handle("/metrics", promhttp.Handler())

	t.logger.New().Infof("server HTTP starting on %s", httpAddr)
	return t.server.ListenAndServe()
}

func (t *HttpTransport) Shutdown(ctx context.Context) error {
	t.logger.New().Infof("stopping HTTP server")

	if t.server != nil {
		return t.server.Shutdown(ctx)
	}
	return nil
}
