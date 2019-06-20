.PHONY: all test build

all:
	drone exec

test: go.mod
	go fmt ./...
	cd db && statik -p mysql -m -Z -f -src=schema/mysql && cd ..
	go test -v --tags="migrations" ./db/...

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/webcron-linux-amd64 cmd/webcron/*.go

go.mod:
	go mod init