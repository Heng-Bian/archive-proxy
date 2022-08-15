package archive

import (
	"archive/zip"
	"errors"
	"github.com/Heng-Bian/archive-proxy/third_party/ranger"
	"io"
	"net/http"
	"net/url"
)

var defaultClient *http.Client = &http.Client{}

func ListZipFiles(zipUrl string, charset string, client *http.Client) (files []string, err error) {
	if client == nil {
		client = defaultClient
	}
	fileNames := make([]string, 0, 10)
	zipReader, err := urlToZipReader(zipUrl, client)
	if err != nil {
		return fileNames, err
	}
	for _, file := range zipReader.File {
		if charset == "" {
			fileNames = append(fileNames, file.Name)
		} else {
			str, err := DecodeString(file.Name, charset)
			if err == nil {
				fileNames = append(fileNames, str)
			} else {
				fileNames = append(fileNames, file.Name)
			}
		}
	}

	return fileNames, nil
}

func UnzipByFileName(zipUrl string, name string, charset string, client *http.Client) (io.ReadCloser, error) {
	if client == nil {
		client = defaultClient
	}
	zipReader, err := urlToZipReader(zipUrl, client)
	if err != nil {
		return nil, err
	}
	for _, file := range zipReader.File {
		fileName := file.Name
		if charset != "" {
			decodeStr, err := DecodeString(fileName, charset)
			if err == nil {
				fileName = decodeStr
			}
		}
		if fileName == name {
			return file.Open()
		}
	}
	return nil, errors.New("file not foud in the zip archive")

}

func UnzipByFileIndex(zipUrl string, index int, client *http.Client) (io.ReadCloser, error) {
	if client == nil {
		client = defaultClient
	}
	zipReader, err := urlToZipReader(zipUrl, client)
	if err != nil {
		return nil, err
	}
	if index > len(zipReader.File) {
		return nil, errors.New("index out of boundary")
	}
	return zipReader.File[index].Open()
}
func urlToReader(httpUrl string, client *http.Client) (*ranger.Reader, error) {
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

func urlToZipReader(zipUrl string, client *http.Client) (*zip.Reader, error) {
	reader, err := urlToReader(zipUrl, client)
	if err != nil {
		return nil, err
	}
	lenth, err := reader.Length()
	if err != nil {
		return nil, err
	}
	zipReader, err := zip.NewReader(reader, lenth)
	if err != nil {
		return nil, err
	}
	return zipReader, nil
}
