package archiveproxy

import (
	"compress/bzip2"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Heng-Bian/archive-proxy/pkg/archive"
	"github.com/ulikunitz/xz"
)

const (
	//parameter name
	targetUrl  = "url"
	charset    = "charset"
	fileIndex  = "index"
	fileFormat = "format"
)

var (
	errReferrer   = errors.New("request does not contain an allowed referrer")
	errDeniedHost = errors.New("request contains a denied host")
	errNotAllowed = errors.New("requested URL is not allowed")
)

type ArchiveStruct struct {
	FileType string
	Files    []string
}

var empty ArchiveStruct

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

	// The Logger used by the archive proxy
	Logger *log.Logger

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
	err := p.allowed(r)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "fail to proxy,err:%s", err)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/healthz") {
		handler = http.HandlerFunc(p.ServeHealthCheck)
	} else if strings.HasPrefix(r.URL.Path, "/list") || strings.HasPrefix(r.URL.Path, "/stream") {
		handler = http.HandlerFunc(p.ServeArchive)
	} else {
		handler = http.HandlerFunc(p.Serve404)
	}
	handler.ServeHTTP(w, r)
}

func (p *Proxy) ServeHealthCheck(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "OK")
}

func (p *Proxy) ServeArchive(w http.ResponseWriter, r *http.Request) {
	targetUrl := r.URL.Query().Get(targetUrl)
	fileFormat := r.URL.Query().Get(fileFormat)
	charset := r.URL.Query().Get(charset)
	index := r.URL.Query().Get(fileIndex)
	if targetUrl == "" {
		w.WriteHeader(500)
		fmt.Fprintf(w, "url must not empty!")
		return
	}
	reader, err := archive.UrlToReader(targetUrl, p.Client)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "fail to crete reader from given url,err:%s", err)
		return
	}
	defer reader.Close()
	if p.IncludeReferer {
		// pass along the referer header from the original request
		copyHeader(reader.Header, r.Header, "referer")
	}
	if len(p.PassRequestHeaders) != 0 {
		copyHeader(reader.Header, r.Header, p.PassRequestHeaders...)
	}
	if fileFormat == "" {
		mimeType, err := archive.DetectMimeTypeThenSeek(reader)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "fail to detect file type,err:%s", err)
			return
		}
		fileFormat = archive.MineTypeTransform(mimeType)
	}

	if strings.HasPrefix(r.URL.Path, "/list") {
		//list archive
		var res ArchiveStruct
		res.FileType = fileFormat
		switch fileFormat {
		case archive.ZIP_TYPE:
			files, err := archive.ListZipFiles(reader, charset)
			res.Files = files
			writeRes(w, res, err)
		case archive.TAR_TYPE:
			files, err := archive.ListTarFiles(reader, charset)
			res.Files = files
			writeRes(w, res, err)
		case archive.RAR_TYPE:
			files, err := archive.ListRarFiles(reader)
			res.Files = files
			writeRes(w, res, err)
		case archive.SEVEN_Z_TYPE:
			files, err := archive.List7zFiles(reader)
			res.Files = files
			writeRes(w, res, err)
		default:
			writeRes(w, res, errors.New("do not support "+res.FileType))
		}

	} else if strings.HasPrefix(r.URL.Path, "/pack") {
		if r.Method != "POST" {
			writeRes(w, empty, errors.New("mehod not allowed"))
		} else {
			var names []string
			err := json.NewDecoder(r.Body).Decode(&names)
			if err != nil {
				writeRes(w, empty, errors.New("mehod not allowed"))
				return
			}
			switch fileFormat {
			case archive.ZIP_TYPE:
				archive.ZipToZip(w, reader, names, charset)
			case archive.TAR_TYPE:
				archive.TarToZip(w, reader, names)
			case archive.SEVEN_Z_TYPE:
				archive.SevenZToZip(w, reader, names)
			case archive.RAR_TYPE:
				archive.RarToZip(w, reader, names)
			default:
				writeRes(w, empty, errors.New("only support zip,tar,7z and rar"))
			}
		}
	} else if strings.HasPrefix(r.URL.Path, "/stream") {
		//return stream
		var isUseFileName bool
		var fileIndex int
		fileName := strings.TrimPrefix(r.URL.Path, "/stream/")

		// file name is empty
		if fileName == r.URL.Path {
			value, err := strconv.Atoi(index)
			if err != nil {
				writeRes(w, empty, err)
				return
			}
			fileIndex = value
			isUseFileName = false
		} else {
			isUseFileName = true
		}
		switch fileFormat {
		case archive.ZIP_TYPE:
			if isUseFileName {
				r, err := archive.UnzipByFileName(reader, fileName, charset)
				writeStream(w, r, err)
			} else {
				r, err := archive.UnzipByFileIndex(reader, fileIndex)
				writeStream(w, r, err)
			}
		case archive.TAR_TYPE:
			if isUseFileName {
				r, err := archive.UnTarByFileName(reader, fileName, charset)
				writeStream(w, r, err)
			} else {
				r, err := archive.UnTarByFileIndex(reader, fileIndex)
				writeStream(w, r, err)
			}
		case archive.RAR_TYPE:
			if isUseFileName {
				r, err := archive.UnRarByFileName(reader, fileName)
				writeStream(w, r, err)
			} else {
				r, err := archive.UnRarByFileIndex(reader, fileIndex)
				writeStream(w, r, err)
			}
		case archive.SEVEN_Z_TYPE:
			if isUseFileName {
				r, err := archive.Un7zByFileName(reader, fileName)
				writeStream(w, r, err)
			} else {
				r, err := archive.Un7zByFileIndex(reader, fileIndex)
				writeStream(w, r, err)
			}
		case archive.GZIP_TYPE:
			r, err := gzip.NewReader(reader)
			writeStream(w, r, err)
		case archive.XZ_TYPE:
			r, err := xz.NewReader(reader)
			writeStream(w, r, err)
		case archive.BZIP2_TYPE:
			r := bzip2.NewReader(reader)
			writeStream(w, r, nil)
		}

	} else {
		w.WriteHeader(404)
	}
}

func (p *Proxy) Serve404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

// allowed determines whether the specified request contains an allowed
// referrer and host.  It returns an error if the request is not
// allowed.
func (p *Proxy) allowed(requst *http.Request) error {
	targetUrl := requst.URL.Query().Get(targetUrl)
	u, err := url.Parse(targetUrl)
	if err != nil {
		return errors.New("invalid target url:" + targetUrl)
	}
	if len(p.AllowHosts) > 0 && !hostMatches(p.AllowHosts, u) {
		return errNotAllowed
	}
	if len(p.DenyHosts) > 0 && hostMatches(p.AllowHosts, u) {
		return errDeniedHost
	}
	if len(p.Referrers) > 0 && !referrerMatches(p.Referrers, requst) {
		return errReferrer
	}
	return nil
}

func writeRes(w http.ResponseWriter, res ArchiveStruct, err error) {
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, err.Error())
		return
	}
	jsonBytes, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, err.Error())
		return
	}
	w.Write(jsonBytes)
}

func writeStream(w http.ResponseWriter, r io.Reader, err error) {
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, err.Error())
		return
	}
	io.Copy(w, r)
}

// hostMatches returns whether the host in u matches one of hosts.
func hostMatches(hosts []string, u *url.URL) bool {
	for _, host := range hosts {
		if u.Hostname() == host {
			return true
		}
		if strings.HasPrefix(host, "*.") && strings.HasSuffix(u.Hostname(), host[2:]) {
			return true
		}
		// Checks whether the host in u is an IP
		if ip := net.ParseIP(u.Hostname()); ip != nil {
			// Checks whether our current host is a CIDR
			if _, ipnet, err := net.ParseCIDR(host); err == nil {
				// Checks if our host contains the IP in u
				if ipnet.Contains(ip) {
					return true
				}
			}
		}
	}

	return false
}

// returns whether the referrer from the request is in the host list.
func referrerMatches(hosts []string, r *http.Request) bool {
	u, err := url.Parse(r.Header.Get("Referer"))
	if err != nil { // malformed or blank header, just deny
		return false
	}
	return hostMatches(hosts, u)
}

// copyHeader copies values for specified headers from src to dst, adding to
// any existing values with the same header name.
func copyHeader(dst, src http.Header, headerNames ...string) {
	for _, name := range headerNames {
		k := http.CanonicalHeaderKey(name)
		for _, v := range src[k] {
			dst.Add(k, v)
		}
	}
}
