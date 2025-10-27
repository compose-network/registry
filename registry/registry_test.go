package registry

import (
	"errors"
	"testing"
)

func TestGetNetworks_FindsHoodi(t *testing.T) {
	r := New()
	nets, err := r.ListNetworks()
	if err != nil {
		t.Fatalf("ListNetworks() error: %v", err)
	}
	if len(nets) == 0 {
		t.Fatalf("expected at least one network")
	}
	var found bool
	for _, n := range nets {
		if n.Slug() == "hoodi" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("did not find 'hoodi' in ListNetworks()")
	}
}

func TestGetNetworkBySlug(t *testing.T) {
	r := New()
	hoodi, err := r.GetNetworkBySlug("hoodi-dev")
	if err != nil {
		t.Fatalf("GetNetworkBySlug() error: %v", err)
	}
	ncfg, err := hoodi.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error: %v", err)
	}
	if ncfg.Name != "hoodi-dev" {
		t.Fatalf("network name = %s, want hoodi", ncfg.Name)
	}
}

func TestGetChains_AndLookups(t *testing.T) {
	r := New()
	hoodi, err := r.GetNetworkBySlug("hoodi-dev")
	if err != nil {
		t.Fatalf("GetNetworkBySlug() error: %v", err)
	}
	chains, err := hoodi.ListChains()
	if err != nil {
		t.Fatalf("ListChains() error: %v", err)
	}
	if len(chains) < 2 {
		t.Fatalf("expected at least 2 chains, got %d", len(chains))
	}
	// Lookup by slug
	a, err := hoodi.GetChainBySlug("rollup-a")
	if err != nil {
		t.Fatalf("GetChainBySlug(rollup-a) error: %v", err)
	}
	acfg, err := a.LoadConfig()
	if err != nil {
		t.Fatalf("chain rollup-a LoadConfig error: %v", err)
	}
	if acfg.ChainID != 77777 {
		t.Fatalf("rollup-a chain_id = %d, want 77777", acfg.ChainID)
	}
	if a.Slug() != "rollup-a" {
		t.Fatalf("slug mismatch: %s", a.Slug())
	}
	// Lookup by id
	b, err := hoodi.GetChainById(88888)
	if err != nil {
		t.Fatalf("GetChainById(88888) error: %v", err)
	}
	bcfg, err := b.LoadConfig()
	if err != nil {
		t.Fatalf("b LoadConfig error: %v", err)
	}
	if bcfg.Name != "rollup-b" {
		t.Fatalf("chain 88888 name = %s, want rollup-b", bcfg.Name)
	}
	if b.Slug() != "rollup-b" {
		t.Fatalf("slug mismatch: %s", b.Slug())
	}
}

func TestGetChainByIdentifier(t *testing.T) {
	r := New()
	c, err := r.GetChainByIdentifier("hoodi-dev/rollup-a")
	if err != nil {
		t.Fatalf("GetChainByIdentifier error: %v", err)
	}
	ccfg, err := c.LoadConfig()
	if err != nil {
		t.Fatalf("c LoadConfig error: %v", err)
	}
	if c.Slug() != "rollup-a" || ccfg.ChainID != 77777 {
		t.Fatalf("unexpected chain/ID: %s (%d)", c.Slug(), ccfg.ChainID)
	}
}

func TestRegistry_ListChains_And_GetChainById(t *testing.T) {
	r := New()
	all, err := r.ListChains()
	if err != nil {
		t.Fatalf("ListChains() error: %v", err)
	}
	if len(all) < 2 {
		t.Fatalf("expected at least 2 chains, got %d", len(all))
	}
	c, err := r.GetChainById(88888)
	if err != nil {
		t.Fatalf("GetChainById error: %v", err)
	}
	cfg, err := c.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if cfg.Name != "rollup-b" {
		t.Fatalf("chain name = %s, want rollup-b", cfg.Name)
	}
}

func TestNewFromDir_Validation(t *testing.T) {
	// Missing networks/ should error
	if _, err := NewFromDir(t.TempDir()); err == nil {
		t.Fatalf("expected error for dir without networks/")
	}
	// Valid: point to embedded data directory on disk
	r, err := NewFromDir("../data")
	if err != nil {
		t.Fatalf("NewFromDir(../data) error: %v", err)
	}
	nets, err := r.ListNetworks()
	if err != nil {
		t.Fatalf("ListNetworks error: %v", err)
	}
	if len(nets) == 0 {
		t.Fatalf("expected at least one network from disk-backed registry")
	}
}

func TestErrors_NotFound(t *testing.T) {
	r := New()
	if _, err := r.GetNetworkBySlug("nope"); err == nil || !errors.Is(err, ErrNetworkNotFound) {
		t.Fatalf("expected ErrNetworkNotFound, got %v", err)
	}
	if _, err := r.GetChainByIdentifier("nope/x"); err == nil || !errors.Is(err, ErrNetworkNotFound) {
		t.Fatalf("expected ErrNetworkNotFound for bad network, got %v", err)
	}
	hoodi, err := r.GetNetworkBySlug("hoodi-dev")
	if err != nil {
		t.Fatalf("GetNetworkBySlug(hoodi): %v", err)
	}
	if _, err := hoodi.GetChainBySlug("x"); err == nil || !errors.Is(err, ErrChainNotFound) {
		t.Fatalf("expected ErrChainNotFound, got %v", err)
	}
}
