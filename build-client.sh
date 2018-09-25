#!/bin/sh

GOOS=linux GOARCH=amd64 go build -o build/linux/idle-tester client.go
GOOS=darwin GOARCH=amd64 go build -o build/mac/idle-tester client.go
GOOS=windows GOARCH=amd64 go build -o build/windows/idle-tester.exe client.go
