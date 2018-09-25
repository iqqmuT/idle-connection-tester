# Idle Connection Tester

Tests if idle TCP connections are timed out.

## Building Client Executable

### Linux 64 bit

    GOOS=linux GOARCH=amd64 go build -o idle-tester client.go

### Windows 64 bit

    GOOS=windows GOARCH=amd64 go build -o idle-tester.exe client.go

### Mac 64 bit

    GOOS=darwin GOARCH=amd64 go build -o idle-tester client.go
