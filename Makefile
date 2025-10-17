# AGPL-3.0
SHELL := bash
NAME := rpcv2-hist
VERSION := 0.1.0
CGO := 0
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: all
all: build

.PHONY: build
build:
	CGO_ENABLED=$(CGO) go build -trimpath -ldflags '$(LDFLAGS)' -o bin/$(NAME) ./cmd/$(NAME)

.PHONY: test
test:
	go test -race -cover ./...

.PHONY: bench
bench:
	go test -bench=. -benchmem ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: image
image:
	docker build -t $(NAME):$(VERSION) .

.PHONY: clean
clean:
	rm -rf bin/

.PHONY: run-local
run-local:
	RPCV2_LOGLEVEL=debug \
	RPCV2_CLICKHOUSE_ADDR=127.0.0.1:9000 \
	RPCV2_CLICKHOUSE_DATABASE=solana \
	go run ./cmd/$(NAME)