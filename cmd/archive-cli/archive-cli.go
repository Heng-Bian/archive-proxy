package main

import (
	"archive/tar"
	"fmt"
	"io"
	"sync"

	"net/http"
	"net/url"

	"github.com/Heng-Bian/archive-proxy/third_party/ranger"
)

func main() {
	targetUrl := "https://aigc-private.obs.cn-north-4.myhuaweicloud.com:443/file_upload/tts-v1.0.2.tar.gz?AccessKeyId=KKTO5OVGE4K81PK6PRGR&Expires=1661681134&Signature=C0zfz4JYwXNjdwKicO01pBa700A%3D"
	url, err := url.Parse(targetUrl)
	if err != nil {
		fmt.Println(err)
	}
	httpRanger := &ranger.HTTPRanger{
		Client: http.DefaultClient,
		URL:    url,
	}
	if err != nil {
		fmt.Println(err)
	}
	r, err := ranger.NewRingBuffReader(httpRanger)
	if err != nil {
		fmt.Println(err)
	}
	tarReader := tar.NewReader(r)
	for {
		header, err := tarReader.Next()
		if err != nil {
			fmt.Println(err)
		}
		if err == io.EOF{
			break
		}
		fmt.Println(header.Name)
	}
	sync.WaitGroup
	go func(){

	}()


}
