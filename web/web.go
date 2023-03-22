package web

import (
	"embed"
)

//go:embed index.html favicon.ico static/*
var EmbedFS embed.FS
