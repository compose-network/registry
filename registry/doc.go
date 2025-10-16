// Package registry exposes embedded Compose network data via
// simple handles (Network, Chain) and on-demand LoadConfig.
//
// Example:
//
//	nets, _ := registry.ListNetworks()
//	hoodi, _ := registry.GetNetworkBySlug("hoodi")
//	hoodi2, _ := registry.GetNetworkById(560048)
//
//	// All chains in this network (handles only)
//	chains, _ := hoodi.ListChains()
//	chainA, _ := hoodi.GetChainBySlug("rollup-a")
//	acfg, _ := chainA.LoadConfig()
//	chainB, _ := hoodi.GetChainById(77777)
//	// Global helpers
//	allChains, _ := registry.ListChains()
//	chainC, _ := registry.GetChainByIdentifier("hoodi/rollup-a")
//	parent := chainC.Network()
package registry
