package main

import (
	"compress/gzip"
	"fmt"
	"io"

	//"io/ioutil"
	"os"

	//"github.com/Heng-Bian/archive-proxy/pkg/archive"
	//"github.com/saracen/go7z"
	"github.com/ulikunitz/xz"
)

func main() {
	
	// r,_:= archive.UrlToReader("http://10.19.32.93:9000/aigc-private/1/test.7z?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=HP7095GYPK1XRT4VI9VP%2F20220819%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20220819T090757Z&X-Amz-Expires=604800&X-Amz-Security-Token=eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJIUDcwOTVHWVBLMVhSVDRWSTlWUCIsImV4cCI6MTY2MDkwMzU4NCwicGFyZW50IjoiYWRtaW4ifQ.OHNByw7zGys1C1e33CNIXV_qIRaHfyxkt_x2aQPKLHDMIbFH6FsySdFMYmARoOGX3IqLaWOj4zxjnoxK5D4fEQ&X-Amz-SignedHeaders=host&versionId=null&X-Amz-Signature=c9786349cc5cfcd9d938e5f1a6473e11c29c65090e5700bbb2ff7881eaf2fb48",nil)
	// lenth,_:=r.Length()
	r,_:= os.Open("C:\\Users\\BIAN\\Desktop\\bian.txt.xz")
	// info,_:=r.Stat()
	// lenth:=info.Size()
	sevenz,_:=xz.NewReader(r)
	sevenz.
	gz,_:=gzip.NewReader()
	gz.Name
	
	

}
