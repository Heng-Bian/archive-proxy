package web

import (
	"embed"
)

//go:embed index.html static/*
var EmbedFS embed.FS
