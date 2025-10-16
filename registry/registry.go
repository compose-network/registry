package registry

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	assets "github.com/compose-network/registry"
)

// Sentinel errors for not-found cases. Use errors.Is to test them.
var (
    ErrNetworkNotFound = errors.New("network not found")
    ErrChainNotFound   = errors.New("chain not found")
)

// Network is a lightweight network handle (slug only).
type Network struct{ slug string }

// Slug returns the network slug.
func (n Network) Slug() string { return n.slug }

// Chain is a lightweight chain handle (slug + parent network).
type Chain struct {
    slug    string
    network Network
}

// Slug returns the chain slug.
func (c Chain) Slug() string     { return c.slug }
// Network returns the parent network handle.
func (c Chain) Network() Network { return c.network }
// Identifier returns "<network>/<slug>".
func (c Chain) Identifier() string { return c.network.slug + "/" + c.slug }

// ChainConfig is decoded from data/networks/<network>/<slug>.toml.
type ChainConfig struct {
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

// NetworkConfig is decoded from data/networks/<slug>/compose.toml.
type NetworkConfig struct {
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

// ListNetworks returns network handles found under data/networks.
func ListNetworks() ([]Network, error) {
	entries, err := assets.FS.ReadDir("data/networks")
	if err != nil {
		return nil, fmt.Errorf("list networks: %w", err)
	}
	slugs := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			slugs = append(slugs, e.Name())
		}
	}
	sort.Strings(slugs)
	out := make([]Network, 0, len(slugs))
	for _, s := range slugs {
		out = append(out, Network{slug: s})
	}
	return out, nil
}

// GetNetworkBySlug returns a handle if data/networks/<slug> exists.
func GetNetworkBySlug(slug string) (Network, error) {
	if _, err := assets.FS.ReadDir(filepath.Join("data/networks", slug)); err != nil {
		return Network{}, fmt.Errorf("%w: %s", ErrNetworkNotFound, slug)
	}
	return Network{slug: slug}, nil
}

// GetNetworkById returns the first network whose L1.ChainID matches.
func GetNetworkById(l1ChainId uint64) (Network, error) {
	nets, err := ListNetworks()
	if err != nil {
		return Network{}, err
	}
	for _, n := range nets {
		cfg, err := n.LoadConfig()
		if err != nil {
			return Network{}, err
		}
		if cfg.L1.ChainID == l1ChainId {
			return n, nil
		}
	}
	return Network{}, ErrNetworkNotFound
}

// ListChains returns chain handles in this network.
func (n Network) ListChains() ([]Chain, error) {
	entries, err := assets.FS.ReadDir(filepath.Join("data/networks", n.slug))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrNetworkNotFound, n.slug)
	}
	slugs := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.EqualFold(name, "compose.toml") || !strings.HasSuffix(name, ".toml") {
			continue
		}
		slugs = append(slugs, strings.TrimSuffix(name, ".toml"))
	}
	sort.Strings(slugs)
	out := make([]Chain, 0, len(slugs))
	for _, s := range slugs {
		out = append(out, Chain{slug: s, network: n})
	}
	return out, nil
}

// GetChainBySlug returns a chain handle if <slug>.toml exists.
func (n Network) GetChainBySlug(slug string) (Chain, error) {
	s := strings.TrimSpace(slug)
	if s == "" {
		return Chain{}, errors.New("empty chain slug")
	}
	// existence check only
	if _, err := assets.FS.ReadFile(filepath.Join("data/networks", n.slug, s+".toml")); err != nil {
		return Chain{}, fmt.Errorf("%w: %s/%s", ErrChainNotFound, n.slug, s)
	}
	return Chain{slug: s, network: n}, nil
}

// GetChainById returns the first chain in this network whose ChainID matches.
func (n Network) GetChainById(l2ChainId uint64) (Chain, error) {
	chains, err := n.ListChains()
	if err != nil {
		return Chain{}, err
	}
	for _, ch := range chains {
		cfg, err := ch.LoadConfig()
		if err != nil {
			return Chain{}, err
		}
		if cfg.ChainID == l2ChainId {
			return ch, nil
		}
	}
	return Chain{}, ErrChainNotFound
}

// ListChains returns all chain handles across all networks.
func ListChains() ([]Chain, error) {
	nets, err := ListNetworks()
	if err != nil {
		return nil, err
	}
	var out []Chain
	for _, n := range nets {
		cs, err := n.ListChains()
		if err != nil {
			return nil, err
		}
		out = append(out, cs...)
	}
	return out, nil
}

// GetChainByIdentifier returns a chain handle for "<network>/<slug>".
func GetChainByIdentifier(identifier string) (Chain, error) {
	s := strings.TrimSpace(identifier)
	if s == "" {
		return Chain{}, errors.New("empty identifier")
	}
	i := strings.IndexByte(s, '/')
	if i <= 0 || i >= len(s)-1 {
		return Chain{}, fmt.Errorf("invalid identifier %q: want <network>/<slug>", identifier)
	}
	netSlug := s[:i]
	chainSlug := s[i+1:]
	n, err := GetNetworkBySlug(netSlug)
	if err != nil {
		return Chain{}, err
	}
	return n.GetChainBySlug(chainSlug)
}

// GetChainById returns the first chain across all networks whose ChainID matches.
func GetChainById(l2ChainId uint64) (Chain, error) {
	nets, err := ListNetworks()
	if err != nil {
		return Chain{}, err
	}
	for _, n := range nets {
		ch, err := n.GetChainById(l2ChainId)
		if err == nil {
			return ch, nil
		}
		if !errors.Is(err, ErrChainNotFound) {
			return Chain{}, err
		}
	}
	return Chain{}, ErrChainNotFound
}

// LoadConfig decodes data/networks/<network>/<slug>.toml for this chain.
func (c Chain) LoadConfig() (ChainConfig, error) {
	s := strings.TrimSpace(c.slug)
	if s == "" {
		return ChainConfig{}, errors.New("empty chain slug")
	}
	path := filepath.Join("data/networks", c.network.slug, s+".toml")
	b, err := assets.FS.ReadFile(path)
	if err != nil {
		return ChainConfig{}, fmt.Errorf("%w: %s/%s", ErrChainNotFound, c.network.slug, s)
	}
	var cfg ChainConfig
	if _, err := toml.Decode(string(b), &cfg); err != nil {
		return ChainConfig{}, fmt.Errorf("decode %s: %w", path, err)
	}
	return cfg, nil
}

// LoadConfig decodes data/networks/<slug>/compose.toml for this network.
func (n Network) LoadConfig() (NetworkConfig, error) {
	b, err := assets.FS.ReadFile(filepath.Join("data/networks", n.slug, "compose.toml"))
	if err != nil {
		return NetworkConfig{}, fmt.Errorf("read compose.toml for %s: %w", n.slug, err)
	}
	var cfg NetworkConfig
	if _, err := toml.Decode(string(b), &cfg); err != nil {
		return NetworkConfig{}, fmt.Errorf("decode compose.toml for %s: %w", n.slug, err)
	}
	return cfg, nil
}
