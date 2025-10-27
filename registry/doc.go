// Package registry exposes embedded Compose network data via
// simple handles (Network, Chain) and on-demand LoadConfig.
//
// Instance usage:
//
//	r := registry.New()
//	nets, _ := r.ListNetworks()
//	hoodi, _ := r.GetNetworkBySlug("hoodi")
//	hoodi2, _ := r.GetNetworkById(560048)
//
//	// All chains in this network (handles only)
//	chains, _ := hoodi.ListChains()
//	chainA, _ := hoodi.GetChainBySlug("rollup-a")
//	acfg, _ := chainA.LoadConfig()
//	chainB, _ := hoodi.GetChainById(77777)
//
//	// Global helpers
//	allChains, _ := r.ListChains()
//	chainC, _ := r.GetChainByIdentifier("hoodi/rollup-a")
//	parent := chainC.Network()
package registry
