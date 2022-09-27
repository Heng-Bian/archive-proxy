package archive

import (
	"errors"
	"github.com/Heng-Bian/httpreader"
	rardecode "github.com/nwaples/rardecode/v2"
	"io"
)

func ListRarFiles(r  *httpreader.Reader) (files []string, err error) {
	fileNames := make([]string, 0, 10)
	rarReader, err := rardecode.NewReader(r)
	if err != nil {
		return nil, err
	}
	for {
		header, err := rarReader.Next()
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

func UnRarByFileName(r *httpreader.Reader, name string) (io.Reader, error) {
	rarReader, err := rardecode.NewReader(r)
	if err != nil {
		return nil, err
	}
	for {
		header, err := rarReader.Next()
		if err != nil {
			//io.EOF is not a error
			if err == io.EOF {
				return nil, errors.New("file not found")
			} else {
				return nil, err
			}
		}
		if name == header.Name {
			return rarReader, nil
		}
	}
}

func UnRarByFileIndex(r *httpreader.Reader, index int) (io.Reader, error) {
	rarReader, err := rardecode.NewReader(r)
	if err != nil {
		return nil, err
	}
	var count int
	for {
		_, err := rarReader.Next()
		if err != nil {
			//io.EOF is not a error
			if err == io.EOF {
				return nil, errors.New("file not found")
			} else {
				return nil, err
			}
		}
		if count == index {
			return rarReader, nil
		}
		count++
	}
}
