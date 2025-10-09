package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	t "github.com/compose-network/registry/internal/types"
)

func main() {
	var in string
	var out string
	flag.StringVar(&in, "in", "data/chainList.toml", "input TOML path")
	flag.StringVar(&out, "out", "generated/chainList.json", "output JSON path")
	flag.Parse()

	var cl t.ChainList
	if _, err := toml.DecodeFile(in, &cl); err != nil {
		fatalf("decode TOML: %v", err)
	}
	cl.Generated = time.Now().UTC().Format(time.RFC3339)
	if cl.Version == "" {
		cl.Version = "0.1.0"
	}

	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		fatalf("mkdir generated: %v", err)
	}
	f, err := os.Create(out)
	if err != nil {
		fatalf("create out: %v", err)
	}
	defer func() { _ = f.Close() }()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cl); err != nil {
		fatalf("encode json: %v", err)
	}
	fmt.Printf("wrote %s (chains=%d)\n", out, len(cl.Chains))
}

func fatalf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(1)
}
