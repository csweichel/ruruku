#!/bin/bash

set -ex
tmpdir=$(mktemp -d)
curl -L -o $tmpdir/protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v3.6.1/protoc-3.6.1-linux-x86_64.zip
mkdir -p $HOME/protoc && cd $HOME/protoc && unzip $tmpdir/protoc.zip
chmod -R 755 $HOME/protoc/include/google
ls -lha $HOME/protoc/bin
rm $tmpdir/protoc.zip
