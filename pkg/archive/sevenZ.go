package archive

import (
	"archive/zip"
	"io"
	"sort"

	"github.com/Heng-Bian/httpreader"
	"github.com/saracen/go7z"
)

func List7zFiles(r *httpreader.Reader) (files []string, err error) {
	fileNames := make([]string, 0, 10)
	length := r.Length
	reader, err := go7z.NewReader(r, length)
	if err != nil {
		return nil, err
	}
	for {
		header, err := reader.Next()
		if err != nil {
			//io.EOF is not a error
			if err == io.EOF {
				return fileNames, nil
			} else {
				return fileNames, err
			}
		}
		fileNames = append(fileNames, header.Name)
	}
}

func Un7zByFileName(r *httpreader.Reader, name string) (io.Reader, error) {
	length := r.Length

	reader, err := go7z.NewReader(r, length)
	if err != nil {
		return nil, err
	}
	for {
		header, err := reader.Next()
		if err != nil {
			//io.EOF is not a error
			if err == io.EOF {
				return nil, ErrFileNotFound
			} else {
				return nil, err
			}
		}
		if header.Name == name {
			return reader, nil
		}
	}
}

func Un7zByFileIndex(r *httpreader.Reader, index int) (io.Reader, error) {
	length := r.Length

	reader, err := go7z.NewReader(r, length)
	if err != nil {
		return nil, err
	}
	var count int
	for {
		_, err := reader.Next()
		if err != nil {
			//io.EOF is not a error
			if err == io.EOF {
				return nil, ErrFileNotFound
			} else {
				return nil, err
			}
		}
		if count == index {
			return reader, nil
		}
		count++
	}
}

func SevenZToZip(w io.Writer, r *httpreader.Reader, names []string) error {
	sevenZ, err := go7z.NewReader(r, r.Length)
	if err != nil {
		return err
	}
	zipWriter := zip.NewWriter(w)
	sort.Strings(names)
	for {
		header, err := sevenZ.Next()
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
				io.Copy(z, sevenZ)
			}
		}
	}
}
