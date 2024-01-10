package main

import (
	"flag"
	"net/http"
	"tolling/common"
)

var httpAddr = flag.String("httpAddr", ":3000", "the listen address of the HTTP server")

func main() {
	flag.Parse()
	logger := common.NewCustomLogger()
	l := logger.New()

	var aggregator Aggregator

	memoryStore := NewMemoryStorer()

	aggregator = NewInvoceAggregator(memoryStore)
	aggregator = NewLogMiddleware(aggregator, logger)

	httpTransport := NewHttpTransport(aggregator, logger)

	http.HandleFunc("/aggregate", httpTransport.handleAggregate)

	l.Infof("server starting on %s", *httpAddr)
	if err := http.ListenAndServe(*httpAddr, nil); err != nil {
		l.Errorf("Failed to start server: \n%v", err)
	}
}
