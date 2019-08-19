.PHONY: all test build

all:
	drone exec

test: go.mod
	go fmt ./...
	go test -v --tags="integration" ./...

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/webcron-linux-amd64 cmd/webcron/*.go

run: build
	build/webcron-linux-amd64

go.mod:
	go mod init