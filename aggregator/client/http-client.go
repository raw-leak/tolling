package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

func (c *AggregatorHttpClient) Aggregate(ctx context.Context, d types.Distance) error {
	l := c.logger.New()
	traceID, ok := ctx.Value(types.KeyTraceID).(string)
	if ok {
		l.WithTraceID(traceID)
	}

	b, err := json.Marshal(d)
	if err != nil {
		l.Error("body marshal failed")
		return err
	}

	req, err := http.NewRequest("POST", c.Endpoint+"/aggregate", bytes.NewReader(b))
	if err != nil {
		l.Errorf("http request failed: \n%v", err)
		return err
	}
	defer req.Body.Close()

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

func (c *AggregatorHttpClient) GetInvoice(ctx context.Context, obuid int) (types.Invoice, error) {
	l := c.logger.New()
	traceID, ok := ctx.Value(types.KeyTraceID).(string)
	if ok {
		l.WithTraceID(traceID)
	}

	req, err := http.NewRequest("GET", c.Endpoint+"/invoice?obuid="+strconv.Itoa(obuid), nil)
	if err != nil {
		l.Errorf("http request failed: \n%v", err)
		return types.Invoice{}, err
	}

	if ok {
		req.Header.Set(types.TraceIDHeader, traceID)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return types.Invoice{}, err
	}

	if res.StatusCode != http.StatusOK {
		l.Error("reading invoice failed")
		return types.Invoice{}, fmt.Errorf("reading invoice failed with [%d] status", res.StatusCode)
	}

	var invoice types.Invoice
	err = json.NewDecoder(res.Body).Decode(&invoice)
	if err != nil {
		l.Errorf("error decoding response body: %v", err)
		return types.Invoice{}, err
	}

	return invoice, nil
}

func (c *AggregatorHttpClient) Shutdown() error {
	c.logger.New().Infof("AggregatorClient is disconnecting")
	return nil
}
