#!/bin/zsh


go get ./...
go mod download
gofmt -w ./..

git describe --tags --abbrev=0 > cmd/hconfig/version.txt

cd cmd/hconfig/
go install
# go build -v

now="$(date)"
printf "build at: %s\n" "$now"