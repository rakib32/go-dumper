#!/bin/bash
set -e

# Building
export GOARCH="amd64"
export GOOS="linux"
export CGO_ENABLED=0

GOPATH=$(go env GOPATH)
REPO_ROOT=$GOPATH/src/github.com/dumper/dumper

go build -o dumper -v $REPO_ROOT/main.go

DOCKER_TAG="1.0.0"

docker build -t dumper/dumper:$DOCKER_TAG .
#docker push [your]/dumper/dumper:$DOCKER_TAG

rm dumper
