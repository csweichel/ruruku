FROM golang:1.11-stretch AS build

# install node
RUN curl -sL https://deb.nodesource.com/setup_11.x | bash - && \
    apt-get install -y nodejs && \
    npm install -g yarn

# install protobuf prerequisites
COPY build/install-protobuf.sh .
RUN apt-get install -y unzip && \
    ./install-protobuf.sh && \
    go get -v google.golang.org/grpc && \
    go get -v github.com/golang/protobuf/protoc-gen-go

# build
COPY . /go/src/github.com/32leaves/ruruku
WORKDIR /go/src/github.com/32leaves/ruruku
RUN cd client && yarn install && cd - && \
    export PATH=$PATH:$HOME/protoc/bin && \
    ./build/protoc.sh && \
    cd client && yarn build && cd - && \
    go get -v ./... && \
    GOXOS=linux GOXARCH=amd64 ./build/build_release.sh

FROM alpine:latest
COPY --from=build /go/src/github.com/32leaves/ruruku/build/release/ruruku_linux_amd64 /app/ruruku
ENTRYPOINT [ "/app/ruruku" ]
CMD ["serve"]
