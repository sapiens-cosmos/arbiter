FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y git build-essential wget unzip nodejs npm

COPY dev/proto/install-golang.sh install-golang.sh

RUN sh install-golang.sh

ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV GOBIN /go/bin
ENV PATH $GOBIN:$GOROOT/bin:$PATH

ENV GOLANG_PROTOBUF_VERSION=1.3.5 \
  GOGO_PROTOBUF_VERSION=1.3.2 \
  GRPC_GATEWAY_VERSION=1.14.7

RUN go get \
  github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v${GRPC_GATEWAY_VERSION} \
  github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@v${GRPC_GATEWAY_VERSION}

RUN go get -u github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

RUN npm install -g swagger-combine

COPY dev/proto/install-protoc.sh install-protoc.sh

RUN sh install-protoc.sh

RUN mkdir /app

WORKDIR /app

# To avoid re-installing dependancies with small changes, perform `go mod download` on build.
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

RUN go install github.com/regen-network/cosmos-proto/protoc-gen-gocosmos
RUN go install github.com/rakyll/statik

VOLUME /app
