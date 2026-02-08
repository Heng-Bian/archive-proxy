package web

import (
	"embed"
)

//go:embed dist/*
var EmbedFS embed.FS
