.PHONY: all test build

all:
	drone exec

test: go.mod
	gotest -v --tags="migrations" ./db/...

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/webcron-linux-amd64 cmd/webcron/*.go

go.mod:
	go mod init