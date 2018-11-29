#!/bin/bash

wd=$(dirname $0)
pushd $wd/.. &> /dev/null

protoc -Ipkg/server/api/ \
       -I../../grpc-ecosystem/grpc-gateway/third_party/googleapis \
       pkg/server/api/v1/*.proto \
       --go_out=plugins=grpc:pkg/server/api/v1

popd  &> /dev/null