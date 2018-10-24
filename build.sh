#! /bin/sh
set -e

if ! [ -x "$(command -v go)" ]; then
    echo "go is not installed"
    exit
fi
if ! [ -x "$(command -v git)" ]; then
    echo "git is not installed"
    exit
fi
if [ -z "$GOPATH" ]; then
    echo "set GOPATH"
    exit
fi

# Building
export GOARCH="amd64"
export GOOS="linux"
export CGO_ENABLED=0
unset GOBIN # for fixing "go get: cannot install cross-compiled binaries when GOBIN is set"

go get .
go build -o dumper -v main.go


docker build . -t dumper:latest
# gcloud docker -- push [url]/dumper/dumper:latest

export GOBIN=$GOPATH/bin # again setting $GOBIN after unset

rm dumper