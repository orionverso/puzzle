#!/bin/sh
# compile.sh binary file.go
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $1 $2

zip $1.zip $1

rm -f $1

