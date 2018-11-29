#!/bin/bash

wd=$(dirname $0)
pushd $wd/.. &> /dev/null
protoc -I pkg/server/api/ pkg/server/api/*.proto --go_out=plugins=grpc:pkg/server/api
popd  &> /dev/null