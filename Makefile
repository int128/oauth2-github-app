.PHONY: all
all:

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	go tool -modfile=tools/go.mod golangci-lint run
