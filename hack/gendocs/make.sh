#!/usr/bin/env bash

pushd $GOPATH/src/pharmer.dev/cloud-controller-manager/hack/gendocs
go run main.go
popd
