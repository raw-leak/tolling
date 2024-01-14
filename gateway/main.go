package main

import (
	"flag"
	"log"
	"tolling/common"
	"tolling/gateway/httpserver"
)

var httpAddr = flag.String("httpAddr", ":6000", "the listen address of the HTTP server")

func main() {
	logger := common.NewCustomLogger()
	targets := map[string]string{"aggregator": "http://localhost:3000"}

	server := httpserver.NewServer(targets, *httpAddr, logger)
	err := server.Start()
	if err != nil {
		log.Fatal("failed to start API-Gateway: ", err)
	}
}
