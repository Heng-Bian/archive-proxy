package archive

import (
	"archive/zip"
	"errors"
	"io"
	"sort"
	"strings"

	"github.com/Heng-Bian/httpreader"
	rardecode "github.com/nwaples/rardecode/v2"
)

func ListRarFiles(r *httpreader.Reader) (files []string, err error) {
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
		if header.IsDir && !strings.HasSuffix(header.Name, "/") {
			fileNames = append(fileNames, header.Name+"/")
		} else {
			fileNames = append(fileNames, header.Name)
		}
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

func RarToZip(w io.Writer, r *httpreader.Reader, names []string) error {
	rarReader, err := rardecode.NewReader(r)
	if err != nil {
		return err
	}
	zipWriter := zip.NewWriter(w)
	sort.Strings(names)
	for {
		header, err := rarReader.Next()
		if err != nil {
			zipWriter.Close()
			//io.EOF is not a error
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		}
		name := header.Name
		if header.IsDir && !strings.HasSuffix(header.Name, "/") {
			name = name + "/"
		}
		if Exists(names, name) {
			z, err := zipWriter.Create(name)
			if err == nil && !strings.HasSuffix(header.Name, "/") {
				io.Copy(z, rarReader)
			}
		}
	}
}
