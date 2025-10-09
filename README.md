<p align="center"><img src="https://framerusercontent.com/images/9FedKxMYLZKR9fxBCYj90z78.png?scale-down-to=512&width=893&height=363" alt="SSV Network"></p>

<img src="https://github.com/ssvlabs/template-repository/actions/workflows/main.yml/badge.svg" alt="Check" />
<a href="https://discord.com/invite/ssvnetworkofficial"><img src="https://img.shields.io/badge/discord-%23ssvlabs-8A2BE2.svg" alt="Discord" /></a>


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
- Public API: `github.com/compose-network/registry/data` (List, Get, Version)

### Layout
- `data/chainList.toml` — human‑authored chain list (name, slug, chain_id, parent, public_rpc, explorer, status, registry_level).
- `data/` — Go package with `//go:embed` and a minimal API.
- `internal/types/` — shared types for dev tools.
- `tools/cmd/{validate,generate}` — offline validator and TOML→JSON generator.

## ⚙️ Build & Dev

Requirements: Go 1.24+

Using the Makefile:
```bash
# Format (goimports) and tidy modules
make format

# Build, test, validate TOML, and generate JSON artifact
make build
make test
make validate
make generate

# Lint (uses the tool declared in go.mod)
make lint
```

Without Makefile:
```bash
go build ./...
go test ./...
go run ./tools/cmd/validate -in data/chainList.toml
go run ./tools/cmd/generate -in data/chainList.toml -out generated/chainList.json
```

## 📦 Usage (as a module)

```bash
go get github.com/compose-network/registry
```

```go
import regdata "github.com/compose-network/registry/data"

chains, _ := regdata.List()      // []Chain{Name, Slug, ChainID, PublicRPC}
ver,    _ := regdata.Version()   // registry version string
```

## 🧪 CI

GitHub Actions runs build, test, validate, lint, and a formatting check on PRs (`.github/workflows/ci.yml`).

## 🤝 Contributing

Issues and PRs are welcome. Please keep the public API minimal and additive. For larger schema changes, open an issue first.

## 📄 License

Repository is distributed under [GPL-3.0](LICENSE).
