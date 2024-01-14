.PHONY: gateway generator receiver calculator aggregator proto 

gateway:
	@go build -o bin/gateway gateway/main.go
	@./bin/gateway

generator:
	@go build -o bin/generator generator/main.go
	@./bin/generator

receiver:
	@go build -o bin/receiver receiver/*.go
	@./bin/receiver

calculator:
	@go build -o bin/calculator calculator/*.go
	@./bin/calculator

aggregator:
	@go build -o bin/aggregator aggregator/*.go
	@./bin/aggregator

# generate gRPC go client
proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative types/ptypes.proto