#!/bin/bash

set -ex
cd /tmp && curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.6.1/protoc-3.6.1-linux-x86_64.zip
mkdir $HOME/protoc && cd $HOME/protoc && unzip /tmp/protoc-3.6.1-linux-x86_64.zip
chmod -R 755 $HOME/protoc/include/google
rm /tmp/protoc-3.6.1-linux-x86_64.zip
