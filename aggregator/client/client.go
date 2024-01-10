package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"tolling/common"
	"tolling/types"
)

type AggregatorHttpClient struct {
	Endpoint string
	logger   common.Logger
}

func NewAggregatorHttpClient(endpoint string, logger common.Logger) *AggregatorHttpClient {
	return &AggregatorHttpClient{Endpoint: endpoint, logger: logger}
}

func (c *AggregatorHttpClient) AggregateInvoice(ctx context.Context, d types.Distance) error {
	l := c.logger.New()
	traceID, ok := ctx.Value(types.KeyTraceID).(string)
	fmt.Println("TRACE ID in client -> ", traceID, ok)
	if ok {
		l.WithTraceID(traceID)
	}

	b, err := json.Marshal(d)
	if err != nil {
		l.Error("body marshal failed")
		return err
	}

	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewReader(b))
	if err != nil {
		l.Errorf("http request failed: \n%v", err)
		return err
	}

	if ok {
		req.Header.Set(types.TraceIDHeader, traceID)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		l.Error("aggregator failed")
		return fmt.Errorf("aggregator failed with [%d] status", res.StatusCode)
	}

	return nil
}
