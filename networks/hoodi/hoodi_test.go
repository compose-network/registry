package hoodi

import "testing"

func TestComposeChains_LoadsRollupsAandB(t *testing.T) {
	chains, err := ComposeChains()
	if err != nil {
		t.Fatalf("ComposeChains() error: %v", err)
	}
	if len(chains) < 2 {
		t.Fatalf("expected at least 2 hoodi chains, got %d", len(chains))
	}
	var a, b *Chain
	for i := range chains {
		switch chains[i].Name {
		case "rollup-a":
			a = &chains[i]
		case "rollup-b":
			b = &chains[i]
		}
	}
	if a == nil || b == nil {
		t.Fatalf("missing rollup-a or rollup-b: a=%v b=%v", a != nil, b != nil)
	}
	if a.ChainID != 77777 {
		t.Fatalf("rollup-a chain_id = %d, want 77777", a.ChainID)
	}
	if b.ChainID != 88888 {
		t.Fatalf("rollup-b chain_id = %d, want 88888", b.ChainID)
	}
}

func TestNetworkConfig_Values(t *testing.T) {
	n, err := NetworkConfig()
	if err != nil {
		t.Fatalf("NetworkConfig() error: %v", err)
	}
	if n.Name != "hoodi" {
		t.Fatalf("network name = %s, want hoodi", n.Name)
	}
	if n.L1.ChainID != 560048 {
		t.Fatalf("l1.chain_id = %d, want 560048", n.L1.ChainID)
	}
	if n.Compose.SP.SuperblockContract != "0x0000000000000000000000000000000000000001" {
		t.Fatalf("compose.sp.superblock_contract = %s, want 0x...01", n.Compose.SP.SuperblockContract)
	}
	if n.Compose.SP.DisputeGameFactory != "0x0000000000000000000000000000000000000002" {
		t.Fatalf("compose.sp.dispute_game_factory = %s, want 0x...02", n.Compose.SP.DisputeGameFactory)
	}
}
