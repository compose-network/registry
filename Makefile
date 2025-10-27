.PHONY: all build test validate generate lint format tidy check-genesis check-generated lint-fix verify

GO ?= go

all: build lint test generate validate check-generated check-genesis

build:
	$(GO) build ./...

test:
	$(GO) test ./...

tidy:
	GOWORK=off GOTOOLCHAIN=go1.23.5 $(GO) mod tidy

validate:
	$(MAKE) -C tools validate

generate:
	$(MAKE) -C tools chainlist-gen

check-genesis:
	$(MAKE) -C tools checkgenesis

lint:
	$(MAKE) -C tools lint

format:
	$(MAKE) -C tools format

lint-fix:
	$(MAKE) -C tools lint-fix

check-generated:
	$(MAKE) -C tools check-generated

verify:
	$(MAKE) -C tools verify
