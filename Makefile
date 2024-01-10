.PHONY: obu receiver calculator aggregator

obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu

receiver:
	@go build -o bin/receiver receiver/*.go
	@./bin/receiver

calculator:
	@go build -o bin/calculator calculator/*.go
	@./bin/calculator

aggregator:
	@go build -o bin/aggregator aggregator/*.go
	@./bin/aggregator



