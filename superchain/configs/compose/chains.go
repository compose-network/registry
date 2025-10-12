package compose

import (
	"embed"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

//go:embed *.toml
var composeTomls embed.FS

// Chain mirrors the subset of OP-style fields we expose to SP.
type Chain struct {
	Name      string `toml:"name"`
	ChainID   uint64 `toml:"chain_id"`
	Addresses struct {
		Mailbox string `toml:"Mailbox"`
	} `toml:"addresses"`
	Compose struct {
		Sequencer struct {
			Host string `toml:"host"`
			Port int    `toml:"port"`
		} `toml:"sequencer"`
	} `toml:"compose"`
}

// ComposeChains decodes all embedded per-chain TOMLs for the compose network.
// It skips the network-level superchain.toml and returns a stable, sorted list.
func ComposeChains() ([]Chain, error) {
	// List *.toml in this directory.
	entries, err := composeTomls.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("list compose configs: %w", err)
	}
	files := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.EqualFold(name, "superchain.toml") || !strings.HasSuffix(name, ".toml") {
			continue
		}
		files = append(files, name)
	}
	sort.Strings(files)

	out := make([]Chain, 0, len(files))
	for _, f := range files {
		raw, err := composeTomls.ReadFile(filepath.Clean(f))
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", f, err)
		}
		var c Chain
		if _, err := toml.Decode(string(raw), &c); err != nil {
			return nil, fmt.Errorf("decode %s: %w", f, err)
		}
		out = append(out, c)
	}
	return out, nil
}
