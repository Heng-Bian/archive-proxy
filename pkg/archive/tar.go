package archive

import (
	"archive/tar"
	"archive/zip"
	"errors"
	"io"
	"sort"

	"github.com/Heng-Bian/httpreader"
)

func ListTarFiles(r *httpreader.Reader, charset string) (files []string, err error) {
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

func UnTarByFileName(r *httpreader.Reader, name string, charset string) (io.Reader, error) {
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

func UnTarByFileIndex(r *httpreader.Reader, index int) (io.Reader, error) {
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

func TarToZip(w io.Writer, r *httpreader.Reader, names []string) error {
	tarReader := tar.NewReader(r)
	zipWriter := zip.NewWriter(w)
	sort.Strings(names)
	for {
		header, err := tarReader.Next()
		if err != nil {
			zipWriter.Close()
			//io.EOF is not a error
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		}
		if Exists(names, header.Name) {
			z, err := zipWriter.Create(header.Name)
			if err == nil {
				io.Copy(z, tarReader)
			}
		}
	}
}
