<p align="center"><img src="https://framerusercontent.com/images/9FedKxMYLZKR9fxBCYj90z78.png?scale-down-to=512&width=893&height=363" alt="SSV Network"></p>

<img src="https://github.com/compose-network/registry/actions/workflows/ci.yml/badge.svg" alt="CI" />
<a href="https://discord.gg/compose-network"><img src="https://img.shields.io/badge/discord-compose--network-5865F2.svg" alt="Discord" /></a>


<p align="center"><b>Compose Registry</b> — Embedded chain list for Compose projects</p>

## ✨ Introduction

Compose Registry is the canonical, machine‑readable source of truth for Compose networks. It ships as a tiny Go module that embeds a curated set of minimal TOML configs so apps, CLIs, and services can depend on one versioned artifact.

By embedding the registry you get:
- Reproducible builds — the exact network catalog travels with your binary.
- Simple runtime selection — choose networks via a flag or config, no external files.
- CI/CD friendly — validate and generate artifacts without reaching out to the network.

The goal is to keep network metadata lightweight, auditable, and easy to consume across the ecosystem.

- Module: `github.com/compose-network/registry`
- Embedded sources: `data/networks/<net>/*.toml` and `data/networks/<net>/compose.toml`
- Public API: `github.com/compose-network/registry/registry`
  - Lightweight handles + on‑demand LoadConfig (instance‑based).

### Layout

- `registry/` — Go package that enumerates `data/networks/*` and exposes instance methods:
  - Listing: `(Registry).ListNetworks()`, `(Registry).ListChains()` (handles; no TOML read)
  - Lookup: `(Registry).GetNetworkBySlug(slug)`, `(Registry).GetNetworkById(l1ChainId)`, `(Registry).GetChainByIdentifier("<network>/<slug>")`, `(Registry).GetChainById(l2ChainId)`
  - Per-network: `Network.LoadConfig()`, `Network.ListChains()`, `Network.GetChainBySlug()`, `Network.GetChainById()`
- `data/` — Data files only (no Go): networks/<net>/*.toml, genesis/, dictionary. Optionally, a generated `chainList.{toml,json}` for external tooling.
- `internal/types/` — shared types for dev tools.
- `tools/cmd/{validate,chainlist-gen}` — validator and generator (configs → chainList.{toml,json}).

#### Schema Notes

- Network slug: the directory name under `data/networks/<network-slug>/`. Used for lookups; must be non‑empty and unique.
- Network name: optional display string in `compose.toml`; display-only, may be empty/non‑unique. Do not use for lookups.
- Chain slug: derived strictly from the filename `<slug>.toml` (TOML cannot override). Used for lookups; must be non‑empty and unique within its network.
- Chain name: optional display string `name` in each `*.toml`; display-only, may be empty/non‑unique. Do not use for lookups.
- Identifier: `<network-slug>/<chain-slug>`; used for cross‑network addressing.

## ⚙️ Build & Dev

Requirements: Go 1.24+

Using the Makefile:
```bash
# Format (goimports) and tidy modules
make format

# Generate optional chainList.{toml,json} (for external tooling), validate, build, test
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
# Optional: generate/validate chainList for external tooling
go run ./tools/cmd/chainlist-gen -base .
go run ./tools/cmd/validate -in data/chainList.toml
```

## 📦 Usage (as a module)

```bash
go get github.com/compose-network/registry
```

```go
import reg "github.com/compose-network/registry/registry"

r := reg.New()
nets, _ := r.ListNetworks()                    // []Network handles
hoodi, _ := r.GetNetworkBySlug("hoodi")        // specific network
fmt.Println(hoodi.Slug())                      // slug (key)
ncfg, _ := hoodi.LoadConfig()                  // display-only fields
fmt.Println(ncfg.Name)

chains, _ := hoodi.ListChains()               // []Chain handles (no config loaded)
chainA, _ := hoodi.GetChainBySlug("rollup-a")
acfg, _ := chainA.LoadConfig()                // fields like ChainID, RPC, etc
chainB, _ := hoodi.GetChainById(77777)

allChains, _ := r.ListChains()                // all L2 chain handles across networks
chain, _ := r.GetChainByIdentifier("hoodi/rollup-a")
ccfg, _ := chain.LoadConfig()
fmt.Println(chain.Slug(), ccfg.Name)          // slug, display name
parent := chain.Network()                     // recover parent Network
```

## API at a Glance

- Constructors
  - New() → Registry — embedded assets (data/)
  - NewFromDir(dir string) (Registry, error) — directory-based data source; dir must contain `networks/`

- Registry methods
  - ListNetworks() → []Network — lists available networks (handles only)
  - GetNetworkBySlug(slug) → Network — handle if networks/<slug> exists
  - GetNetworkById(l1ChainId) → Network — scan via LoadConfig()
  - ListChains() → []Chain — lists all chains across all networks (handles only)
  - GetChainByIdentifier("<network>/<slug>") → Chain — resolves identifier
  - GetChainById(l2ChainId) → Chain — scan via Network.GetChainById()

- Network methods
  - Slug() string — unique network slug
  - LoadConfig() → NetworkConfig — loads compose.toml when needed
  - ListChains() → []Chain — lists chain handles in this network
  - GetChainBySlug(slug) → Chain — returns a chain handle if <slug>.toml exists
  - GetChainById(l2ChainId) → Chain — scan via Chain.LoadConfig()

- Chain methods
  - Slug() string — unique chain slug
  - Network() Network — parent network handle
  - Identifier() string — "<network>/<slug>"
  - LoadConfig() → ChainConfig — loads <slug>.toml when needed

### Error Contract

When a network or chain is not found, functions return typed sentinel errors:
- ErrNetworkNotFound
- ErrChainNotFound

You can test with errors.Is:

```go
if errors.Is(err, reg.ErrNetworkNotFound) { /* handle */ }
```

## 🧪 CI

GitHub Actions runs build, test, validate, lint, and a formatting check on PRs (`.github/workflows/ci.yml`).

## 🤝 Contributing

Issues and PRs are welcome. Please keep the public API minimal and additive. For larger schema changes, open an issue first.

## 📄 License

Repository is distributed under [GPL-3.0](LICENSE).

## Credit

Inspired by Optimism's [SuperChain Registry](https://github.com/ethereum-optimism/superchain-registry/)
