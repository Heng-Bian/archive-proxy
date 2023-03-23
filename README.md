# archive-proxy

![GitHub](https://img.shields.io/github/license/Heng-Bian/archive-proxy)
![GitHub](https://img.shields.io/badge/build-pass-green)

archive-proxy is a archive proxy server written in go. It features:
 - list all archive items for the given archive url (zip, tar, rar, 7z)
 - autodetect the file type
 - random access to the single item of big archive on the url (eg. s3 url)
 - easy to build and deploy, since it's pure go
 - support multiple compressed file, eg. zip, tar, rar, 7z, gz, xz, bzip2

I use the archive-proxy to list the archive and download the chosen item of it
before I download the entire archive on the network. It's very useful for big zip
file.

## List the archive items

GET /list

### request parameter

|name|location|type|required|description|
|---|---|---|---|---|
|url|query|string| YES |the archive URL|
|charset|query|string| NO |specify the charset name, default utf-8|
|format|query|string| NO |indicate the file format, autodetect by default|

### request example
```
GET /list?url=https://golang.google.cn/dl/go1.20.1.windows-amd64.zip HTTP/1.1
Host: localhost:8080
```

### response example (Not showing all)
File ending with "/" means directory

```json
{
    "FileType": "zip",
    "Files": [
        "go/",
        "go/CONTRIBUTING.md",
        "go/LICENSE",
        "go/PATENTS",
        "go/README.md",
        "go/SECURITY.md",
        "go/VERSION",
        "go/api/",
        "go/api/README"
    ]
}
```

## Download a single item

GET /stream/{entry}

### request parameter

|name|location|type|required|description|
|---|---|---|---|---|
|entry|path|string| YES |entry name in the Files array. |
|url|query|string| YES |the archive URL|
|charset|query|string| NO |specify the charset name, default utf-8|
|format|query|string| NO |indicate the file format, autodetect by default|

### Request example
```
GET /stream/go/README.md?url=https://golang.google.cn/dl/go1.20.1.windows-amd64.zip HTTP/1.1
Host: localhost:8080
```

### Response example
binary file stream

## Download mutiple item to a zip file

POST /pack

|name|location|type|required|description|
|---|---|---|---|---|
|url|query|string| YES |the archive URL|
|charset|query|string| NO |specify the charset name, default utf-8|
|format|query|string| NO |indicate the file format, autodetect by default|
|body|body|array[string]| YES |entry name array|

### request example
```
POST /pack?url=https://golang.google.cn/dl/go1.20.1.windows-amd64.zip HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Content-Length: 134

[
    "go/CONTRIBUTING.md",
    "go/LICENSE",
    "go/PATENTS",
    "go/README.md",
    "go/SECURITY.md",
    "go/api/README"
]
```
### response example

zip binary stream

## Quick start

### Build
```
git clone https://github.com/Heng-Bian/archive-proxy.git
cd archive-proxy/cmd/archive-server
go build
```
###  Run
For help info  
`./archive-server -help`  
Start the service with port 8080  
`./archive-server -port 8080`

### Access in a browser
After runing the archive-server,
visit `http://localhost:8080`

## Mechanism
archiver-proxy offers an random access to archive item before download the entire
file. archiver-proxy itself do not cache any data and erverything is based on stream. The archive file on the network MUST support HTTP Range request. Fortunately, the common server such as nginx and Minio support it.

For how the Reader implements io.ReaderAt, io.Reader, and io.Seeker depending on HTTP Range Requests, see `https://github.com/Heng-Bian/httpreader`. It's the cleanest and most efficient implementation.

## Warning
Decompressing is a complex topic. archiver-proxy directly exposed to the open Internet is extremely vulnerable. Some archive(eg. zipbomb)may be evil and result in infinite loop or large bandwidth usage. It's recommended that the archive-proxy is deployed on the cloud(eg. k8s) with limited resource.

It's DevOps duty to protect the archive-proxy from untrusted user.

However, issues and PRs are always welcome :)
