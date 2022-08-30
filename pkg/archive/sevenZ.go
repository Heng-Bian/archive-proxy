package archive

import (
	"io"

	"github.com/Heng-Bian/archive-proxy/third_party/ranger"
	"github.com/saracen/go7z"
)

func List7zFiles(r *ranger.Reader) (files []string, err error) {
	fileNames := make([]string, 0, 10)
	length:= r.Length
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

func Un7zByFileName(r *ranger.Reader, name string) (io.Reader, error) {
	length:= r.Length

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

func Un7zByFileIndex(r *ranger.Reader, index int) (io.Reader, error) {
	length:= r.Length

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
