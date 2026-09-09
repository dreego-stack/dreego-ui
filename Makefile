.PHONY: test test-race vet check example

PORT ?= 8091

test:
	go test ./...
	go test ./_tests

test-race:
	go test -race ./...
	go test -race ./_tests

vet:
	go vet ./...
	go vet ./_tests

check: test vet

example:
	PORT=$(PORT) go run ./example
