#!/bin/bash

wd=$(dirname $0)
pushd $wd/.. &> /dev/null

protoc -Ipkg/api/ \
       pkg/api/v1/*.proto \
       --go_out=plugins=grpc:pkg/api

[ ! -d client/src/api ] && mkdir client/src/api
protoc \
    --plugin="protoc-gen-ts=client/node_modules/.bin/protoc-gen-ts" \
    --js_out="import_style=commonjs,binary:client/src/api" \
    --ts_out=service=true:"client/src/api" \
    -Ipkg/api/ \
    pkg/api/v1/*.proto

popd  &> /dev/null