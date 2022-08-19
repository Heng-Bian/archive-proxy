package archive

import (
	"errors"
	"io"
	"github.com/gabriel-vasile/mimetype"
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
