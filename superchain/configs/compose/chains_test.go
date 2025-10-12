package compose

import "testing"

func TestComposeChains_LoadsRollupsAandB(t *testing.T) {
	chains, err := ComposeChains()
	if err != nil {
		t.Fatalf("ComposeChains() error: %v", err)
	}
	if len(chains) < 2 {
		t.Fatalf("expected at least 2 compose chains, got %d", len(chains))
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
		t.Fatalf("missing rollup-a or rollup-b in ComposeChains() output: a=%v b=%v", a != nil, b != nil)
	}

	// rollup-a expectations
	if a.ChainID != 77777 {
		t.Fatalf("rollup-a chain_id = %d, want 77777", a.ChainID)
	}
	if a.Addresses.Mailbox == "" {
		t.Fatalf("rollup-a addresses.Mailbox empty")
	}
	if a.Compose.Sequencer.Host != "optimism-stack-geth" || a.Compose.Sequencer.Port != 9898 {
		t.Fatalf("rollup-a sequencer = %s:%d, want optimism-stack-geth:9898", a.Compose.Sequencer.Host, a.Compose.Sequencer.Port)
	}

	// rollup-b expectations
	if b.ChainID != 88888 {
		t.Fatalf("rollup-b chain_id = %d, want 88888", b.ChainID)
	}
	if b.Addresses.Mailbox == "" {
		t.Fatalf("rollup-b addresses.Mailbox empty")
	}
	if b.Compose.Sequencer.Host != "optimism-stack-2-geth" || b.Compose.Sequencer.Port != 9898 {
		t.Fatalf("rollup-b sequencer = %s:%d, want optimism-stack-2-geth:9898", b.Compose.Sequencer.Host, b.Compose.Sequencer.Port)
	}
}
