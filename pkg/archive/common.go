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

	ErrFileNotFound  = errors.New("file not found in archive")
	ErrOutOfBoundary = errors.New("file index out of archive boundary")
)

const (
	RAR_MIME_TYPE     = "application/x-rar-compressed"
	ZIP_MIME_TYPE     = "application/zip"
	TAR_MIME_TYPE     = "application/x-tar"
	SEVEN_Z_MIME_TYPE = "application/x-7z-compressed"

	GZIP_MIME_TYPE  = "application/x-gzip"
	BZIP2_MIME_TYPE = "application/x-bzip2"
	XZ_MIME_TYPE    = "application/x-xz"
	DEFALUT_MIME    = "application/octet-stream"

	RAR_TYPE     = "rar"
	ZIP_TYPE     = "zip"
	TAR_TYPE     = "tar"
	SEVEN_Z_TYPE = "7z"

	GZIP_TYPE  = "gzip"
	BZIP2_TYPE = "bzip2"
	XZ_TYPE    = "xz"
)

func ListSupprotedFileFormat() []string {
	supprot := make([]string, 0, 7)
	supprot = append(supprot, RAR_TYPE)
	supprot = append(supprot, ZIP_TYPE)
	supprot = append(supprot, TAR_TYPE)
	supprot = append(supprot, SEVEN_Z_TYPE)

	supprot = append(supprot, GZIP_TYPE)
	supprot = append(supprot, BZIP2_TYPE)
	supprot = append(supprot, XZ_TYPE)
	return supprot
}

func MineTypeTransform(mimeType string) string {
	switch mimeType {
	case RAR_MIME_TYPE:
		return RAR_TYPE
	case ZIP_MIME_TYPE:
		return ZIP_TYPE
	case TAR_MIME_TYPE:
		return TAR_TYPE
	case SEVEN_Z_MIME_TYPE:
		return SEVEN_Z_TYPE
	case GZIP_MIME_TYPE:
		return GZIP_TYPE
	case BZIP2_MIME_TYPE:
		return BZIP2_TYPE
	case XZ_MIME_TYPE:
		return XZ_TYPE
	case DEFALUT_MIME:
		return ""
	}
	return ""
}

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
