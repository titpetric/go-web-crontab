#!/bin/bash
set -e
PROJECT=$(<.project)

docker run --rm -v $(pwd):/go/src/github.com/$PROJECT -w /go/src/github.com/$PROJECT -e GOOS=${OS} -e GOARCH=${ARCH} -e CGO_ENABLED=0 -e GOARM=7 titpetric/golang dep ensure -update -v
