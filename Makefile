.PHONY: all
all:

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	$(MAKE) -C tools
	./tools/bin/golangci-lint run
