package web
import(
	"embed"
)

//go:embed static/*
var EmbedFS embed.FS