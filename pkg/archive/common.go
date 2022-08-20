package archive

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/Heng-Bian/archive-proxy/third_party/ranger"
	"github.com/gabriel-vasile/mimetype"
	"golang.org/x/text/encoding/ianaindex"
)
var (
	defaultClient = &http.Client{}

	FILE_NOT_FOUND = errors.New("file not found in archive")
	INDEX_OUT_BOUDARY = errors.New("file index out of archive bounday")
)
	

const (
	RAR_MIME_TYPE  = "application/x-rar-compressed"
	ZIP_MIME_TYPE  = "application/zip"
	GZIP_MIME_TYPE = "application/x-gzip"
	TAR_MIME_TYPE     = "application/x-tar"
	SEVEN_Z_MIME_TYPE = "application/x-7z-compressed"
	BZIP2_MIME_TYPE   = "application/x-bzip2"
	XZ_MIME_TYPE = "application/x-xz"
	DEFALUT_MIME = "application/octet-stream"
)

func DetectMimeTypeThenSeek(r io.Reader) (string, error) {
	mime, err := mimetype.DetectReader(r)
	if err != nil {
		return "", err
	}
	seeker, ok := r.(io.Seeker)
	if ok {
		seeker.Seek(0, io.SeekStart)
	} else {
		return "", errors.New("fail to seek after detecting mime type")
	}
	return mime.String(), nil
}

func UrlToReader(httpUrl string, client *http.Client) (*ranger.Reader, error) {
	if client == nil {
		client = defaultClient
	}
	url, err := url.Parse(httpUrl)
	if err != nil {
		return nil, err
	}
	httpRanger := &ranger.HTTPRanger{
		Client: client,
		URL:    url,
	}
	reader, err := ranger.NewReader(httpRanger)
	if err != nil {
		return nil, err
	}
	return reader, err
}

func DecodeString(src string, name string) (string, error) {
	encoding, err := ianaindex.IANA.Encoding(name)
	if err != nil {
		return "", err
	}
	target, err := encoding.NewDecoder().String(src)
	if err != nil {
		return "", err
	}
	return target, nil
}