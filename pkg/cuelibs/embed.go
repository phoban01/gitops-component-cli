package cuelibs

import "embed"

// Files contains files to include
//
//go:embed cue.mod
//go:embed ocm
var Files embed.FS
