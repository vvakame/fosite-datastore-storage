#!/bin/sh -eux

cd `dirname $0`

# build tools
rm -rf build-cmd/
mkdir build-cmd

export GOBIN=`pwd -P`/build-cmd
go install golang.org/x/tools/cmd/goimports
go install golang.org/x/lint/golint
go install honnef.co/go/tools/cmd/staticcheck
go install github.com/google/wire/cmd/wire
