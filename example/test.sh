#!/bin/bash -eux

cd `dirname $0`

targets=`find . -type f \( -name '*.go' -and -not -iwholename '*vendor*'  -and -not -iwholename '*node_modules*' \)`
packages=`go list ./...`

export PATH=$(pwd)/build-cmd:$PATH
which goimports golint staticcheck wire
goimports -w $targets
for package in $packages
do
    go vet $package
done
# golint -set_exit_status -min_confidence 0.6 $packages
staticcheck $packages
go generate $packages

go test $packages -p 1 -coverpkg=`go list -m`/... -covermode=atomic -coverprofile=coverage.txt $@
