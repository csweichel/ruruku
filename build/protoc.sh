#!/bin/bash

wd=$(dirname $0)
pushd $wd/.. &> /dev/null

protoc -Ipkg/server/api/ \
       pkg/server/api/v1/*.proto \
       --go_out=plugins=grpc:pkg/server/api

popd  &> /dev/null