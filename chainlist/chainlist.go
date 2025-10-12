package chainlist

import (
	"strconv"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	assets "github.com/compose-network/registry"
	t "github.com/compose-network/registry/internal/types"
)

// Entry mirrors the public schema and adds a derived Slug.
type Entry struct {
	Name                 string   `json:"name"`
	Identifier           string   `json:"identifier"`
	Slug                 string   `json:"slug"`
	ChainID              uint64   `json:"chainId"`
	RPC                  []string `json:"rpc"`
	Explorers            []string `json:"explorers"`
	DataAvailabilityType string   `json:"dataAvailabilityType"`
	Parent               Parent   `json:"parent"`
	GasPayingToken       string   `json:"gasPayingToken,omitempty"`
	FaultProofs          *Faults  `json:"faultProofs,omitempty"`
}

type Parent struct {
	Type  string `json:"type"`
	Chain string `json:"chain"`
}

type Faults struct {
	Status string `json:"status"`
}

var (
	parsed struct {
		once sync.Once
		cl   t.ChainListTOML
		err  error
	}
)

// List returns all chains from the embedded chainList.toml.
func List() ([]Entry, error) {
	if err := ensureParsed(); err != nil {
		return nil, err
	}
	out := make([]Entry, 0, len(parsed.cl.Chains))
	for _, c := range parsed.cl.Chains {
		e := Entry{
			Name:                 c.Name,
			Identifier:           c.Identifier,
			Slug:                 deriveSlug(c.Identifier),
			ChainID:              c.ChainID,
			RPC:                  append([]string(nil), c.RPC...),
			Explorers:            append([]string(nil), c.Explorers...),
			DataAvailabilityType: c.DataAvailabilityType,
			Parent:               Parent{Type: c.Parent.Type, Chain: c.Parent.Chain},
			GasPayingToken:       c.GasPayingToken,
		}
		if c.FaultProofs != nil {
			e.FaultProofs = &Faults{Status: c.FaultProofs.Status}
		}
		out = append(out, e)
	}
	return out, nil
}

// Get returns a chain by slug (derived from identifier suffix).
func Get(slug string) (Entry, bool, error) {
	if err := ensureParsed(); err != nil {
		return Entry{}, false, err
	}
	for _, c := range parsed.cl.Chains {
		if deriveSlug(c.Identifier) == slug {
			e := Entry{
				Name:                 c.Name,
				Identifier:           c.Identifier,
				Slug:                 slug,
				ChainID:              c.ChainID,
				RPC:                  append([]string(nil), c.RPC...),
				Explorers:            append([]string(nil), c.Explorers...),
				DataAvailabilityType: c.DataAvailabilityType,
				Parent:               Parent{Type: c.Parent.Type, Chain: c.Parent.Chain},
				GasPayingToken:       c.GasPayingToken,
			}
			if c.FaultProofs != nil {
				e.FaultProofs = &Faults{Status: c.FaultProofs.Status}
			}
			return e, true, nil
		}
	}
	return Entry{}, false, nil
}

// Version returns a synthetic version string based on count.
func Version() (string, error) {
	if err := ensureParsed(); err != nil {
		return "", err
	}
	return "vchains-" + strconv.Itoa(len(parsed.cl.Chains)), nil
}

func ensureParsed() error {
	parsed.once.Do(func() {
		b, err := assets.FS.ReadFile("data/chainList.toml")
		if err != nil {
			parsed.err = err
			return
		}
		parsed.err = toml.Unmarshal(b, &parsed.cl)
	})
	return parsed.err
}

func deriveSlug(identifier string) string {
	if idx := strings.LastIndex(identifier, "/"); idx >= 0 && idx < len(identifier)-1 {
		return identifier[idx+1:]
	}
	return identifier
}

// GetByIdentifier returns an entry by its full identifier (e.g., "hoodi/rollup-a").
func GetByIdentifier(identifier string) (Entry, bool, error) {
	if err := ensureParsed(); err != nil {
		return Entry{}, false, err
	}
	for _, c := range parsed.cl.Chains {
		if c.Identifier == identifier {
			e := Entry{
				Name:                 c.Name,
				Identifier:           c.Identifier,
				Slug:                 deriveSlug(c.Identifier),
				ChainID:              c.ChainID,
				RPC:                  append([]string(nil), c.RPC...),
				Explorers:            append([]string(nil), c.Explorers...),
				DataAvailabilityType: c.DataAvailabilityType,
				Parent:               Parent{Type: c.Parent.Type, Chain: c.Parent.Chain},
				GasPayingToken:       c.GasPayingToken,
			}
			if c.FaultProofs != nil {
				e.FaultProofs = &Faults{Status: c.FaultProofs.Status}
			}
			return e, true, nil
		}
	}
	return Entry{}, false, nil
}

// ListByNetwork returns all entries that belong to a given L1 parent network (e.g., "hoodi").
func ListByNetwork(network string) ([]Entry, error) {
	if err := ensureParsed(); err != nil {
		return nil, err
	}
	out := make([]Entry, 0)
	for _, c := range parsed.cl.Chains {
		if c.Parent.Chain == network {
			e := Entry{
				Name:                 c.Name,
				Identifier:           c.Identifier,
				Slug:                 deriveSlug(c.Identifier),
				ChainID:              c.ChainID,
				RPC:                  append([]string(nil), c.RPC...),
				Explorers:            append([]string(nil), c.Explorers...),
				DataAvailabilityType: c.DataAvailabilityType,
				Parent:               Parent{Type: c.Parent.Type, Chain: c.Parent.Chain},
				GasPayingToken:       c.GasPayingToken,
			}
			if c.FaultProofs != nil {
				e.FaultProofs = &Faults{Status: c.FaultProofs.Status}
			}
			out = append(out, e)
		}
	}
	return out, nil
}
