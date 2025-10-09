package data

import (
	"sync"

	"github.com/BurntSushi/toml"

	_ "embed"
)

//go:embed chainList.toml
var chainListBytes []byte

// Chain is a minimal exported view of a chain entry.
type Chain struct {
	Name      string
	Slug      string
	ChainID   int64
	PublicRPC string
}

type chainList struct {
	Version   string  `toml:"version"`
	Generated string  `toml:"generated"`
	Chains    []chain `toml:"chain"`
}

type chain struct {
	Name             string `toml:"name"`
	Slug             string `toml:"slug"`
	ChainID          int64  `toml:"chain_id"`
	Parent           string `toml:"parent"`
	PublicRPC        string `toml:"public_rpc"`
	Explorer         string `toml:"explorer"`
	DataAvailability string `toml:"data_availability_type"`
	Status           string `toml:"status"`
	RegistryLevel    int    `toml:"registry_level"`
}

var (
	parsed struct {
		once sync.Once
		cl   chainList
		err  error
	}
)

// List returns all chains from the embedded chainList.toml.
func List() ([]Chain, error) {
	if err := ensureParsed(); err != nil {
		return nil, err
	}
	out := make([]Chain, 0, len(parsed.cl.Chains))
	for _, c := range parsed.cl.Chains {
		out = append(out, Chain{Name: c.Name, Slug: c.Slug, ChainID: c.ChainID, PublicRPC: c.PublicRPC})
	}
	return out, nil
}

// Get returns a chain by slug.
func Get(slug string) (Chain, bool, error) {
	if err := ensureParsed(); err != nil {
		return Chain{}, false, err
	}
	for _, c := range parsed.cl.Chains {
		if c.Slug == slug {
			return Chain{Name: c.Name, Slug: c.Slug, ChainID: c.ChainID, PublicRPC: c.PublicRPC}, true, nil
		}
	}
	return Chain{}, false, nil
}

// Version returns the top-level version string.
func Version() (string, error) {
	if err := ensureParsed(); err != nil {
		return "", err
	}
	return parsed.cl.Version, nil
}

func ensureParsed() error {
	parsed.once.Do(func() {
		parsed.err = toml.Unmarshal(chainListBytes, &parsed.cl)
	})
	return parsed.err
}
