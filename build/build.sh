#!/bin/bash
export GOOS=linux
go build ../cmd/archive-server/archive-server.go
docker build -t bianheng/archive
