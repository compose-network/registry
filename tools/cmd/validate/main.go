package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	t "github.com/compose-network/registry/internal/types"
)

var caip2 = regexp.MustCompile(`^eip155:\d+$`)

func main() {
	var in string
	flag.StringVar(&in, "in", "data/chainList.toml", "input TOML path")
	flag.Parse()
	var cl t.ChainList
	if _, err := toml.DecodeFile(in, &cl); err != nil {
		fatalf("decode TOML: %v", err)
	}
	if err := validate(cl); err != nil {
		fatalf("validation failed: %v", err)
	}
	fmt.Println("validation ok")
}

func validate(cl t.ChainList) error {
	seenName := map[string]bool{}
	seenSlug := map[string]bool{}
	seenID := map[int64]bool{}
	for i, c := range cl.Chains {
		if c.Name == "" || c.Slug == "" {
			return fmt.Errorf("chain[%d]: name/slug required", i)
		}
		if seenName[c.Name] {
			return fmt.Errorf("duplicate name: %s", c.Name)
		}
		if seenSlug[c.Slug] {
			return fmt.Errorf("duplicate slug: %s", c.Slug)
		}
		if c.ChainID <= 0 {
			return fmt.Errorf("chain[%d]: invalid chain_id", i)
		}
		if seenID[c.ChainID] {
			return fmt.Errorf("duplicate chain_id: %d", c.ChainID)
		}
		seenName[c.Name], seenSlug[c.Slug], seenID[c.ChainID] = true, true, true

		if c.Parent != "" && !caip2.MatchString(c.Parent) {
			return fmt.Errorf("chain[%d]: parent must be CAIP-2 eip155:<id>", i)
		}
		if err := mustURL(c.PublicRPC); err != nil {
			return fmt.Errorf("chain[%d] public_rpc: %w", i, err)
		}
		if c.Explorer != "" {
			if err := mustURL(c.Explorer); err != nil {
				return fmt.Errorf("chain[%d] explorer: %w", i, err)
			}
		}
		switch strings.ToLower(c.Status) {
		case "active", "incubating", "deprecated":
		default:
			return fmt.Errorf("chain[%d]: invalid status %q", i, c.Status)
		}
		if c.RegistryLevel < 0 || c.RegistryLevel > 3 {
			return fmt.Errorf("chain[%d]: registry_level out of range", i)
		}
	}
	return nil
}

func mustURL(s string) error {
	if s == "" {
		return errors.New("empty URL")
	}
	u, err := url.Parse(s)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("invalid url: %s", s)
	}
	return nil
}

func fatalf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(1)
}
