<p align="center"><img src="https://framerusercontent.com/images/9FedKxMYLZKR9fxBCYj90z78.png?scale-down-to=512&width=893&height=363" alt="SSV Network"></p>

<img src="https://github.com/compose-network/registry/actions/workflows/ci.yml/badge.svg" alt="CI" />
<a href="https://discord.gg/compose-network"><img src="https://img.shields.io/badge/discord-compose--network-5865F2.svg" alt="Discord" /></a>


<p align="center"><b>Compose Registry</b> — Embedded chain list for Compose projects</p>

## ✨ Introduction

Compose Registry is the canonical, machine‑readable source of truth for Compose networks. It ships as a tiny Go module that embeds a curated list of chains (TOML) so apps, CLIs, and services can depend on one versioned artifact.

By embedding the registry you get:
- Reproducible builds — the exact network catalog travels with your binary.
- Simple runtime selection — choose networks via a flag or config, no external files.
- CI/CD friendly — validate and generate artifacts without reaching out to the network.

The goal is to keep network metadata lightweight, auditable, and easy to consume across the ecosystem.

- Module: `github.com/compose-network/registry`
- Embedded source: `data/chainList.toml`
- Public API: `github.com/compose-network/registry/chainlist`
  - List, Get, GetByIdentifier, ListByNetwork, Version

### Layout

- `data/chainList.toml` — auto-generated chain summary file:
  - `[[chains]]` entries with fields: `name`, `identifier` (e.g., `hoodi/rollup-a`), `chain_id`, `rpc`[], `explorers`[] and nested `parent.{type, chain}`.
  - We derive a short `slug` at runtime from `identifier`’s suffix (e.g., `rollup-a`).
- `chainlist/` — Go package that embeds `data/chainList.toml` and exposes List/Get/GetByIdentifier/ListByNetwork/Version.
- `networks/<network>/` — Go package that embeds `data/networks/<network>/*.toml` and `compose.toml` and exposes ComposeChains/NetworkConfig.
- `data/` — Data files only (no Go): chainList.{toml,json}, networks/<net>/*.toml, genesis/, dictionary.
- `internal/types/` — shared types for dev tools.
- `tools/cmd/{validate,chainlist-gen}` — validator and generator (configs → chainList.{toml,json}).

## ⚙️ Build & Dev

Requirements: Go 1.24+

Using the Makefile:
```bash
# Format (goimports) and tidy modules
make format

# Generate from configs, validate, build, test
make generate
make validate
make build
make test

# Lint (uses the tool declared in go.mod)
make lint
```

Without Makefile (generate both files from configs):
```bash
go build ./...
go test ./...
go run ./tools/cmd/chainlist-gen -base .
go run ./tools/cmd/validate -in data/chainList.toml
```

## 📦 Usage (as a module)

```bash
go get github.com/compose-network/registry
```

```go
import "github.com/compose-network/registry/chainlist"

chains, _ := chainlist.List()      // []Entry with arrays: RPC, Explorers
one, ok, _ := chainlist.Get("rollup-a")
byID, ok2, _ := chainlist.GetByIdentifier("hoodi/rollup-a")
hoodiOnly, _ := chainlist.ListByNetwork("hoodi")
ver, _ := chainlist.Version()      // e.g., vchains-2
```

## 🧪 CI

GitHub Actions runs build, test, validate, lint, and a formatting check on PRs (`.github/workflows/ci.yml`).

## 🤝 Contributing

Issues and PRs are welcome. Please keep the public API minimal and additive. For larger schema changes, open an issue first.

## 📄 License

Repository is distributed under [GPL-3.0](LICENSE).

## Credit

Inspired by Optimism's [SuperChain Registry](https://github.com/ethereum-optimism/superchain-registry/)
