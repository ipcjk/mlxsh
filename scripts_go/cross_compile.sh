#!/bin/bash

# X-compile everything ;-)
env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/mlxsh *.go
env GOOS=windows GOARCH=amd64 go  build -o bin/mlxsh.exe *.go
env GOOS=darwin GOARCH=amd64 go  build -ldflags="-s -w"  -o bin/mlxsh.mac *.go


