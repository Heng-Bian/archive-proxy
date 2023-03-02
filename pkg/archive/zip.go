package archive

import (
	"archive/zip"
	"errors"
	"github.com/Heng-Bian/httpreader"
	"io"
	"sort"
)

func ListZipFiles(r *httpreader.Reader, charset string) (files []string, err error) {
	fileNames := make([]string, 0, 10)
	lenth := r.Length
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
	lenth := r.Length
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
	lenth := r.Length
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

func ZipToZip(w io.Writer, r *httpreader.Reader, names []string, charset string) error {
	zipReader, err := zip.NewReader(r, r.Length)
	if err != nil {
		return err
	}
	zipWriter := zip.NewWriter(w)
	sort.Strings(names)
	for _, file := range zipReader.File {
		fileName := file.Name
		if charset != "" {
			decodeStr, err := DecodeString(fileName, charset)
			if err == nil {
				fileName = decodeStr
			}
		}
		if Exists(names, fileName) {
			zr, err := file.Open()
			if err == nil {
				zw, err := zipWriter.Create(fileName)
				if err == nil {
					io.Copy(zw, zr)
				}
			}
		}
	}
	zipWriter.Close()
	return nil
}
