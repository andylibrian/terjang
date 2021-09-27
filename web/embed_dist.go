package web

import "embed"

//go:embed dist/*
var WebDistFs embed.FS
