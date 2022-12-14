package archive

import (
	"archive/zip"
	"errors"
	"io"
	"github.com/Heng-Bian/httpreader"
)

func ListZipFiles(r *httpreader.Reader, charset string) (files []string, err error) {
	fileNames := make([]string, 0, 10)
	lenth:= r.Length
	zipReader, err := zip.NewReader(r, lenth)
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

func UnzipByFileName(r *httpreader.Reader, name string, charset string) (io.Reader, error) {
	lenth:= r.Length
	zipReader, err := zip.NewReader(r, lenth)
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

func UnzipByFileIndex(r *httpreader.Reader, index int) (io.Reader, error) {
	lenth:= r.Length
	zipReader, err := zip.NewReader(r, lenth)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if index > len(zipReader.File) {
		return nil, ErrOutOfBoundary
	}
	return zipReader.File[index].Open()
}
