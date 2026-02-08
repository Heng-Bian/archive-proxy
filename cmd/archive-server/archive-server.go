package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Heng-Bian/archive-proxy/internal/archiveproxy"
	"github.com/Heng-Bian/archive-proxy/web"
)

var (
	port               = flag.String("port", "8080", "port to listen on")
	ip                 = flag.String("ip", "0.0.0.0", "address to listen on")
	allowHosts         = flag.String("allowHosts", "", "comma separated list of allowed remote hosts")
	denyHosts          = flag.String("denyHosts", "", "comma separated list of denied remote hosts")
	referrers          = flag.String("referrers", "", "comma separated list of allowed referring hosts")
	includeReferer     = flag.Bool("includeReferer", true, "include referer header in remote requests")
	passRequestHeaders = flag.String("passRequestHeaders", "", "comma separatetd list of request headers to pass to remote server")
)

func main() {
	parse("ARCHIVE")
	flag.Parse()
	log.SetFlags(log.Llongfile | log.LUTC)
	proxy := archiveproxy.NewProxy(http.DefaultClient)
	if *allowHosts != "" {
		proxy.AllowHosts = strings.Split(*allowHosts, ",")
	}
	if *denyHosts != "" {
		proxy.DenyHosts = strings.Split(*denyHosts, ",")
	}
	if *referrers != "" {
		proxy.Referrers = strings.Split(*referrers, ",")
	}
	proxy.IncludeReferer = *includeReferer
	if *passRequestHeaders != "" {
		proxy.PassRequestHeaders = strings.Split(*passRequestHeaders, ",")
	}
	addr := *ip + ":" + *port
	server := &http.Server{
		Addr: addr,
	}
	// Serve the React app from the dist subdirectory
	distFS, _ := fs.Sub(web.EmbedFS, "dist")
	http.Handle("/", http.FileServer(http.FS(distFS)))
	http.Handle("/healthz", http.HandlerFunc(proxy.ServeHealthCheck))
	http.Handle("/list", http.HandlerFunc(proxy.ServeArchive))
	http.Handle("/pack", http.HandlerFunc(proxy.ServeArchive))
	http.Handle("/stream", http.HandlerFunc(proxy.ServeArchive))
	server.ListenAndServe()
}

func parse(p string) {
	update(p, flag.CommandLine)
}

func update(p string, fs *flag.FlagSet) {
	// Build a map of explicitly set flags.
	set := map[string]interface{}{}
	fs.Visit(func(f *flag.Flag) {
		set[f.Name] = nil
	})

	fs.VisitAll(func(f *flag.Flag) {
		envVar := fmt.Sprintf("%s_%s", p, strings.ToUpper(f.Name))
		envVar = strings.Replace(envVar, "-", "_", -1)
		if val := os.Getenv(envVar); val != "" {
			if _, defined := set[f.Name]; !defined {
				fs.Set(f.Name, val)
			}
		}
		f.Usage = fmt.Sprintf("%s [%s]", f.Usage, envVar)
	})
}
