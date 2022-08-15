package archive

import (
	"archive/tar"
	"errors"
	"io"
	"net/http"
)

func ListTarFiles(tarUrl string, charset string, client *http.Client) (files []string, err error) {
	fileNames := make([]string, 0, 10)
	if client == nil {
		client = defaultClient
	}
	reader, err := urlToReader(tarUrl, client)
	if err != nil {
		return fileNames, err
	}
	tarReader := tar.NewReader(reader)
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

func UnTarByFileName(tarUrl string, name string, charset string, client *http.Client) (io.Reader, error) {
	if client == nil {
		client = defaultClient
	}
	reader, err := urlToReader(tarUrl, client)
	if err != nil {
		return nil, err
	}
	tarReader := tar.NewReader(reader)
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

func UnTarByFileIndex(tarUrl string, index int, client *http.Client) (io.Reader, error) {
	if client == nil {
		client = defaultClient
	}
	reader, err := urlToReader(tarUrl, client)
	if err != nil {
		return nil, err
	}
	tarReader := tar.NewReader(reader)
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
