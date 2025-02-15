.PHONY: all
all:

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	go tool github.com/golangci/golangci-lint/cmd/golangci-lint run
