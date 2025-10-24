.PHONY: tidy lint format chainlist-gen validate checkgenesis lint-fix check-generated verify

GO ?= go
TOOLCHAIN ?= go1.24.9
export GOWORK=off
export GOTOOLCHAIN=$(TOOLCHAIN)

# Defaults for paths (caller can override)
ROOT ?= ..
BASE ?= ..
# Outputs/inputs are relative to BASE inside the tools; the generator/validator
# will join them with -base. Keep these paths BASE-relative (no BASE prefix)
# to avoid double-joining like "../../data/..." when invoked via -base.
OUT_TOML ?= data/chainList.toml
OUT_JSON ?= data/chainList.json
IN ?= data/chainList.toml

tidy:
	$(GO) mod tidy

lint: tidy
	$(GO) run github.com/golangci/golangci-lint/v2/cmd/golangci-lint run $(ROOT)/...

format: tidy
	$(GO) run golang.org/x/tools/cmd/goimports -w $(ROOT)

chainlist-gen: tidy
	$(GO) run ./cmd/chainlist-gen -base $(BASE) -out-toml $(OUT_TOML) -out-json $(OUT_JSON)

validate: tidy
	$(GO) run ./cmd/validate -in $(BASE)/$(IN)

checkgenesis: tidy
	$(GO) run ./cmd/checkgenesis -base $(BASE)

lint-fix: tidy
	$(MAKE) format
	$(GO) run github.com/golangci/golangci-lint/v2/cmd/golangci-lint run --fix $(ROOT)/...

check-generated: chainlist-gen
	@git -C $(ROOT) diff --quiet -- data/chainList.toml data/chainList.json || (echo 'error: chainList files are stale; run make generate and commit' && git -C $(ROOT) --no-pager diff -- data/chainList.toml data/chainList.json && exit 1)

verify: chainlist-gen validate lint check-generated
