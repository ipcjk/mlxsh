#!/bin/bash
env GOOS=linux GOARCH=amd64 go build -o bin/brocadecli.linux *.go
env GOOS=windows GOARCH=amd64 go  build -o bin/brocadecli.exe *.go
env GOOS=darwin GOARCH=amd64 go  build -o bin/brocadecli.mac *.go


