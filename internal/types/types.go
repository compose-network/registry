package types

type ChainList struct {
	Version   string  `toml:"version" json:"version"`
	Generated string  `toml:"generated" json:"generated"`
	Chains    []Chain `toml:"chain" json:"chains"`
}

type Chain struct {
	Name             string `toml:"name" json:"name"`
	Slug             string `toml:"slug" json:"slug"`
	ChainID          int64  `toml:"chain_id" json:"chain_id"`
	Parent           string `toml:"parent" json:"parent,omitempty"`
	PublicRPC        string `toml:"public_rpc" json:"public_rpc"`
	Explorer         string `toml:"explorer" json:"explorer,omitempty"`
	DataAvailability string `toml:"data_availability_type" json:"data_availability_type,omitempty"`
	Status           string `toml:"status" json:"status"`
	RegistryLevel    int    `toml:"registry_level" json:"registry_level"`
}
