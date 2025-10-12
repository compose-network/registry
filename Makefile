.PHONY: all build test validate generate lint format tidy checkgenesis check-generated

GO ?= go

all: generate validate build lint test

build:
	$(GO) build ./...

test:
	$(GO) test ./...

validate:
	$(GO) run ./tools/cmd/validate -in data/chainList.toml

generate:
	$(GO) run ./tools/cmd/chainlist-gen -base . -out-toml data/chainList.toml -out-json data/chainList.json

checkgenesis:
	$(GO) run ./tools/cmd/checkgenesis

lint:
	$(GO) tool golangci-lint run ./...

format:
	$(GO) tool goimports -w .

tidy:
	$(GO) mod tidy

check-generated:
	$(GO) run ./tools/cmd/chainlist-gen -base . -out-toml data/chainList.toml -out-json data/chainList.json
	@git diff --quiet -- data/chainList.toml data/chainList.json || (echo 'error: chainList files are stale; run make generate and commit' && git --no-pager diff -- data/chainList.toml data/chainList.json && exit 1)
