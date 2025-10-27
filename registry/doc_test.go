package registry

import (
	"fmt"
)

func ExampleNew() {
	r := New()
	nets, _ := r.ListNetworks()
	found := false
	for _, n := range nets {
		if n.Slug() == "hoodi" {
			found = true
			break
		}
	}
	fmt.Println(found)
	// Output: true
}

func ExampleRegistry_ListNetworks() {
	r := New()
	nets, _ := r.ListNetworks()
	found := false
	for _, n := range nets {
		if n.Slug() == "hoodi" {
			found = true
			break
		}
	}
	if found {
		fmt.Println("ok")
	} else {
		fmt.Println("missing")
	}
	// Output: ok
}

func ExampleNetwork_ListChains() {
	r := New()
	n, _ := r.GetNetworkBySlug("hoodi")
	chains, _ := n.ListChains()
	found := false
	for _, c := range chains {
		if c.Slug() == "rollup-a" {
			found = true
			break
		}
	}
	if found {
		fmt.Println("ok")
	} else {
		fmt.Println("missing")
	}
	// Output: ok
}

func ExampleRegistry_GetChainByIdentifier() {
	r := New()
	c, _ := r.GetChainByIdentifier("hoodi/rollup-a")
	fmt.Println(c.Identifier())
	// Output: hoodi/rollup-a
}
