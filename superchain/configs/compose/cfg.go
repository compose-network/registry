package compose

import (
	_ "embed"
	"fmt"

	"github.com/BurntSushi/toml"
)

//go:embed superchain.toml
var superchainToml []byte

// Network models the network-level compose superchain config we expose.
type Network struct {
	Name                 string         `toml:"name"`
	L1                   L1             `toml:"l1"`
	ProtocolVersionsAddr string         `toml:"protocol_versions_addr"`
	Compose              ComposeSection `toml:"compose"`
}

type L1 struct {
	ChainID   uint64 `toml:"chain_id"`
	PublicRPC string `toml:"public_rpc"`
	Explorer  string `toml:"explorer"`
}

type ComposeSection struct {
	SP SP `toml:"sp"`
}

type SP struct {
	SuperblockContract string `toml:"superblock_contract"`
	DisputeGameFactory string `toml:"dispute_game_factory"`
}

// NetworkConfig decodes the embedded superchain.toml and returns it.
func NetworkConfig() (Network, error) {
	var n Network
	if _, err := toml.Decode(string(superchainToml), &n); err != nil {
		return Network{}, fmt.Errorf("decode compose superchain.toml: %w", err)
	}
	return n, nil
}
