package web

import "embed"

// WebClientAssets storing some assets into the binary with the go embedded file system
//go:embed client/*
var WebClientAssets embed.FS
