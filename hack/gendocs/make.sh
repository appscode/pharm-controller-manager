#!/usr/bin/env bash

pushd $GOPATH/src/github.com/pharmer/cloud-controller-manager/hack/gendocs
go run main.go
popd
