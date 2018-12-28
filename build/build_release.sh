#!/bin/bash

VERSION="_"
if [ ! -z "$TRAVIS_TAG" ]; then
    VERSION="_${TRAVIS_TAG}_"
fi

GOXOS=${GOXOS:="linux darwin windows"}
GOXARCH=${GOXARCH:="386 amd64"}

go get -v github.com/mitchellh/gox
go get -v github.com/GeertJohan/go.rice/rice
go get -v github.com/inconshreveable/mousetrap
go get -v github.com/konsorten/go-windows-terminal-sequences
mkdir build/release
rice embed-go -v -i github.com/32leaves/ruruku/pkg/server
CGO_ENABLED=0 gox -os="$GOXOS" -arch="386 amd64" -output="build/release/ruruku$VERSION{{.OS}}_{{.Arch}}" -ldflags "-s -w -X github.com/32leaves/ruruku/pkg/version.Tag=`git describe --abbrev=0` -X github.com/32leaves/ruruku/pkg/version.Rev=`git rev-parse --short HEAD` -X 'github.com/32leaves/ruruku/pkg/version.BuildDate=`date`'" -verbose github.com/32leaves/ruruku
