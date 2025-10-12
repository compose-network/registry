// Package hoodi provides access to the embedded Compose network
// configuration for the Hoodi L1 environment. It loads per-chain
// TOMLs (rollup-*.toml) and the network-level compose.toml from
// the embedded data/ tree.
//
// Primary entry points:
//   - ComposeChains: returns decoded per-chain configs for Hoodi
//   - NetworkConfig: returns the network-level L1 + Compose SP config
package hoodi
