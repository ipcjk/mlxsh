#!/bin/bash
env GOOS=linux GOARCH=amd64 go build -o bin/mlxsh *.go
#env GOOS=windows GOARCH=amd64 go  build -o bin/mlxsh.exe *.go
#env GOOS=darwin GOARCH=amd64 go  build -o bin/mlxsh.mac *.go


