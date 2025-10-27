package assets

import "embed"

// FS contains the embedded registry data files.
//
//go:embed data/**
var FS embed.FS
