package hoodi

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	assets "github.com/compose-network/registry"
)

type Chain struct {
	Name      string `toml:"name"`
	ChainID   uint64 `toml:"chain_id"`
	PublicRPC string `toml:"public_rpc"`
	Explorer  string `toml:"explorer"`
	Addresses struct {
		Mailbox string `toml:"Mailbox"`
	} `toml:"addresses"`
	Genesis struct {
		L2Time uint64 `toml:"l2_time"`
	} `toml:"genesis"`
	Compose struct {
		Sequencer struct {
			Host string `toml:"host"`
			Port int    `toml:"port"`
		} `toml:"sequencer"`
	} `toml:"compose"`
}

// ComposeChains decodes all embedded per-chain TOMLs for the hoodi network.
func ComposeChains() ([]Chain, error) {
	entries, err := assets.FS.ReadDir("data/networks/hoodi")
	if err != nil {
		return nil, fmt.Errorf("list hoodi configs: %w", err)
	}
	files := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.EqualFold(name, "compose.toml") || !strings.HasSuffix(name, ".toml") {
			continue
		}
		files = append(files, name)
	}
	sort.Strings(files)

	out := make([]Chain, 0, len(files))
	for _, f := range files {
		raw, err := assets.FS.ReadFile(filepath.Join("data/networks/hoodi", f))
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

// Network represents compose (network-level) config for hoodi.
type Network struct {
	Name string `toml:"name"`
	L1   struct {
		ChainID   uint64 `toml:"chain_id"`
		PublicRPC string `toml:"public_rpc"`
		Explorer  string `toml:"explorer"`
	} `toml:"l1"`
	Compose struct {
		SP struct {
			SuperblockContract string `toml:"superblock_contract"`
			DisputeGameFactory string `toml:"dispute_game_factory"`
		} `toml:"sp"`
	} `toml:"compose"`
}

func NetworkConfig() (Network, error) {
	var n Network
	b, err := assets.FS.ReadFile("data/networks/hoodi/compose.toml")
	if err != nil {
		return Network{}, fmt.Errorf("read compose.toml: %w", err)
	}
	if _, err := toml.Decode(string(b), &n); err != nil {
		return Network{}, fmt.Errorf("decode compose.toml: %w", err)
	}
	return n, nil
}
