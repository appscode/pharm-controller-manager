#!/usr/bin/env bash

pushd $GOPATH/src/github.com/appscode/pharm-controller-manager/hack/gendocs
go run main.go
popd
