package main

import (
	"log"
	"net/http"
	"github.com/Heng-Bian/archive-proxy/internal/archiveproxy"
)

func main() {
	log.SetFlags(log.Llongfile | log.LUTC)
	proxy := archiveproxy.NewProxy(http.DefaultClient)
	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: proxy,
	}
	server.ListenAndServe()
}
