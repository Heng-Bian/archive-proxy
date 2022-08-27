package archive

import (
	"archive/tar"
	"errors"
	"io"

	"github.com/Heng-Bian/archive-proxy/third_party/ranger"
)

func ListTarFiles(r  *ranger.RingBuffReader, charset string) (files []string, err error) {
	fileNames := make([]string, 0, 10)
	tarReader := tar.NewReader(r)
	for {
		header, err := tarReader.Next()
		if err != nil {
			//io.EOF is not a error
			if err == io.EOF {
				return fileNames, nil
			} else {
				return fileNames, err
			}
		}
		entryName := header.Name
		if charset != "" {
			str, err := DecodeString(entryName, charset)
			if err == nil {
				entryName = str
			}
		}
		fileNames = append(fileNames, entryName)
	}
}

func UnTarByFileName(r *ranger.RingBuffReader, name string, charset string) (io.Reader, error) {
	tarReader := tar.NewReader(r)
	for {
		header, err := tarReader.Next()
		if err != nil {
			//io.EOF is not a error
			if err == io.EOF {
				return nil, errors.New("file not found")
			} else {
				return nil, err
			}
		}
		entryName := header.Name
		if charset != "" {
			str, err := DecodeString(entryName, charset)
			if err == nil {
				entryName = str
			}
		}
		if name == entryName {
			return tarReader, nil
		}
	}
}

func UnTarByFileIndex(r *ranger.RingBuffReader, index int) (io.Reader, error) {
	tarReader := tar.NewReader(r)
	var count int
	for {
		_, err := tarReader.Next()
		if err != nil {
			//io.EOF is not a error
			if err == io.EOF {
				return nil, errors.New("file not found")
			} else {
				return nil, err
			}
		}
		if count == index {
			return tarReader, nil
		}
		count++
	}
}
