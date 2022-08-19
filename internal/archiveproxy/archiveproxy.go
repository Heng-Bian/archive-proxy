package archiveproxy

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	zip = "zip"
	tar = "tar"
)

type ArchiveStruct struct {
	Archive_type string
	Files []string
}

type Proxy struct {
	// client used to fetch remote URLs
	Client *http.Client

	// AllowHosts specifies a list of remote hosts that archives can be
	// proxied from.  An empty list means all hosts are allowed.
	AllowHosts []string

	// DenyHosts specifies a list of remote hosts that archives cannot be
	// proxied from.
	DenyHosts []string

	// Referrers, when given, requires that requests to the archive
	// proxy come from a referring host. An empty list means all
	// hosts are allowed.
	Referrers []string

	// IncludeReferer controls whether the original Referer request header
	// is included in remote requests.
	IncludeReferer bool

	// FollowRedirects controls whether archiveproxy will follow redirects or not.
	FollowRedirects bool

	// The Logger used by the archive proxy
	Logger *log.Logger

	// Timeout specifies a time limit for requests served by this Proxy.
	// If a call runs for longer than its time limit, a 504 Gateway Timeout
	// response is returned.  A Timeout of zero means no timeout.
	Timeout time.Duration

	// The User-Agent used by archiveproxy when requesting origin archive
	UserAgent string

	// PassRequestHeaders identifies HTTP headers to pass from inbound
	// requests to the proxied server.
	PassRequestHeaders []string
}

func NewProxy(client *http.Client) *Proxy {
	proxy := new(Proxy)
	proxy.Client = client
	return proxy
}

// ServeHTTP handles incoming requests.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler
	if r.URL.Path == "/favicon.ico" {
		return // ignore favicon requests
	}
	
	switch r.URL.Path {
	case "/healthz":
		handler = http.HandlerFunc(p.ServeHealthCheck)
	case "/zip":
		handler = http.HandlerFunc(p.Serve404)
	case "/tar":
		handler = http.HandlerFunc(p.ServeZip)
	default:
		handler = http.HandlerFunc(p.Serve404)
	}
	handler.ServeHTTP(w, r)
}

func (p *Proxy) ServeHealthCheck(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "OK")
}

func (p *Proxy) Serve404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

// /zip/${innerPath}?url=https://example.com/example.zip
// url=https://example.com/examle.zip
func (p *Proxy) ServeAuto(w http.ResponseWriter, r *http.Request){
	http.DetectContentType()
}
func (p *Proxy) ServeZip(w http.ResponseWriter, r *http.Request) {

}

func (p *Proxy) ServeTar(w http.ResponseWriter, r *http.Request) {

}

func (p *Proxy) ServeGzip(w http.ResponseWriter, r *http.Request) {

}

func (p *Proxy) Serve7z(w http.ResponseWriter, r *http.Request) {

}

func (p *Proxy) ServeBzip2(w http.ResponseWriter, r *http.Request) {

}

func (p *Proxy) ServeXz(w http.ResponseWriter, r *http.Request) {

}

func (p *Proxy) ServeRar(w http.ResponseWriter, r *http.Request) {

}

// security check before handle requests
func (p *Proxy) PreCheck(w http.ResponseWriter, r *http.Request){
	//TODO
	
}
