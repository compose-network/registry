.PHONY: all build test validate generate lint format tidy

GO ?= go

all: lint build test validate generate

build:
	$(GO) build ./...

test:
	$(GO) test ./...

validate:
	$(GO) run ./tools/cmd/validate -in data/chainList.toml

generate:
	$(GO) run ./tools/cmd/generate -in data/chainList.toml -out generated/chainList.json

lint:
	$(GO) tool golangci-lint run ./...

format:
	$(GO) tool goimports -w .

tidy:
	$(GO) mod tidy
