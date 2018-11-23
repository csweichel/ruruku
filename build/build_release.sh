#!/bin/bash

VERSION="_"
if [ ! -z "$TRAVIS_TAG" ]; then
    VERSION="_${TRAVIS_TAG}_"
fi

go get -v github.com/mitchellh/gox
go get -v github.com/GeertJohan/go.rice/rice
go get -v github.com/inconshreveable/mousetrap
go get -v github.com/konsorten/go-windows-terminal-sequences
mkdir build/release
rice embed-go -v -i github.com/32leaves/ruruku/pkg/rurukulib/server
gox -os="linux darwin windows" -arch="386 amd64" -output="build/release/ruruku$VERSION{{.OS}}_{{.Arch}}" -ldflags "-s -w -X cmd.Rev=`git rev-parse --short HEAD`" -verbose github.com/32leaves/ruruku
