package web

import "embed"

//go:embed client/*
var WebClientAssets embed.FS
