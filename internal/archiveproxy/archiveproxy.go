package archiveproxy

import (
	"compress/bzip2"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Heng-Bian/archive-proxy/pkg/archive"
	"github.com/Heng-Bian/archive-proxy/third_party/ranger"
	"github.com/ulikunitz/xz"
)

const (
	//parameter name
	targetUrl  = "url"
	charset    = "charset"
	fileIndex  = "index"
	fileFormat = "format"
)

type ArchiveStruct struct {
	FileType string
	Files    []string
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

	if strings.HasPrefix(r.URL.Path, "/healthz") {

		handler = http.HandlerFunc(p.ServeHealthCheck)

	} else if strings.HasPrefix(r.URL.Path, "/list") {

		handler = http.HandlerFunc(p.ServeArchive)

	} else if strings.HasPrefix(r.URL.Path, "/stream") {

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

	var reader *ranger.Reader

	if targetUrl == "" {
		fmt.Fprintf(w, "url must not empty!")
		w.WriteHeader(500)
		return
	} else {
		r, err := archive.UrlToReader(targetUrl, p.Client)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "fail to crete reader from given url,err:%s", err)
			return
		} else {
			reader = r
		}
	}
	if fileFormat == "" {
		mimeType, err := archive.DetectMimeTypeThenSeek(reader)
		if err != nil {
			fmt.Fprintf(w, "fail to detect file type,err:%s", err)
			w.WriteHeader(500)
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

	} else if strings.HasPrefix(r.URL.Path, "/stream") {
		//return stream
		var isUseFileName bool
		var fileIndex int
		fileName := strings.TrimPrefix(r.URL.Path, "/stream/")

		// file name is empty
		if fileName == r.URL.Path {
			value, err := strconv.Atoi(index)
			if err != nil {
				var empty ArchiveStruct
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

// security check before handle requests
func (p *Proxy) PreCheck(w http.ResponseWriter, r *http.Request) {
	//TODO

}

func writeRes(w http.ResponseWriter, res ArchiveStruct, err error) {
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, err.Error())
		return
	}
	jsonBytes, err := json.Marshal(res)
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
