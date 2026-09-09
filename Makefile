.PHONY: test test-race vet check

test:
	go test ./...

test-race:
	go test -race ./...

vet:
	go vet ./...

check: test vet
