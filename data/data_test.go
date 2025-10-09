package data

import "testing"

func TestListAndGet(t *testing.T) {
	chains, err := List()
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(chains) < 2 {
		t.Fatalf("expected at least 2 chains, got %d", len(chains))
	}

	// Ensure rollup-a exists with expected ID
	found := false
	for _, c := range chains {
		if c.Slug == "rollup-a" {
			found = true
			if c.ChainID != 77777 {
				t.Fatalf("rollup-a chain_id = %d, want 77777", c.ChainID)
			}
		}
	}
	if !found {
		t.Fatalf("rollup-a not found in List()")
	}

	if got, ok, err := Get("rollup-b"); err != nil || !ok {
		t.Fatalf("Get(rollup-b) = (%v,%v,%v), want ok=true and no error", got, ok, err)
	} else if got.ChainID != 88888 {
		t.Fatalf("rollup-b chain_id = %d, want 88888", got.ChainID)
	}
}

func TestVersion(t *testing.T) {
	v, err := Version()
	if err != nil {
		t.Fatalf("Version() error: %v", err)
	}
	if v == "" {
		t.Fatalf("Version() empty, want non-empty")
	}
}
