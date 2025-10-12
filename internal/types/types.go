package types

type ChainListTOML struct {
	Chains []ChainListEntry `toml:"chains"`
}

type ChainListEntry struct {
	Name                 string               `toml:"name" json:"name"`
	Identifier           string               `toml:"identifier" json:"identifier"`
	ChainID              uint64               `toml:"chain_id" json:"chainId"`
	RPC                  []string             `toml:"rpc" json:"rpc"`
	Explorers            []string             `toml:"explorers" json:"explorers"`
	DataAvailabilityType string               `toml:"data_availability_type" json:"dataAvailabilityType"`
	Parent               ChainListEntryParent `toml:"parent" json:"parent"`
	GasPayingToken       string               `toml:"gas_paying_token,omitempty" json:"gasPayingToken,omitempty"`
	FaultProofs          *FaultProofs         `toml:"fault_proofs,omitempty" json:"faultProofs,omitempty"`
}

type ChainListEntryParent struct {
	Type  string `toml:"type" json:"type"`
	Chain string `toml:"chain" json:"chain"`
}

type FaultProofs struct {
	Status string `toml:"status" json:"status"`
}
