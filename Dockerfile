FROM golang:1.11-alpine AS build
RUN apk add protobuf bash curl git nodejs nodejs-npm && \
     npm install -g yarn && \
     go get -v google.golang.org/grpc && \
     go get -v github.com/golang/protobuf/protoc-gen-go
COPY . /go/src/github.com/32leaves/ruruku
WORKDIR /go/src/github.com/32leaves/ruruku
RUN cd client && yarn install && yarn build; cd - && \
    ./build/protoc.sh && \
    go get -v ./... && \
    GOXOS=linux GOXARCH=amd64 ./build/build_release.sh

FROM alpine:latest
COPY --from=build /go/src/github.com/32leaves/ruruku/build/release/ruruku_linux_amd64 /app/ruruku
ENTRYPOINT [ "/app/ruruku" ]
CMD ["serve"]
