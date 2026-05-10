BINARY ?= bin/supacrawl

.PHONY: build test vet validate run fmt clean

build:
	mkdir -p $(dir $(BINARY))
	go build -o $(BINARY) ./cmd/supacrawl

test:
	go test ./...

vet:
	go vet ./...

validate:
	./scripts/validate-local.sh

run:
	go run ./cmd/supacrawl $(ARGS)

fmt:
	gofmt -w cmd internal

clean:
	rm -rf bin dist
