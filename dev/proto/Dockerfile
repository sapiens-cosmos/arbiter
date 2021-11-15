FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y git build-essential wget unzip nodejs npm

RUN wget https://dl.google.com/go/go1.17.3.linux-amd64.tar.gz
RUN tar -xvf go1.17.3.linux-amd64.tar.gz
RUN mv go /usr/local

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

RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v3.19.1/protoc-3.19.1-linux-x86_64.zip
RUN unzip protoc-3.19.1-linux-x86_64.zip -d /protoc
RUN cp -r /protoc/bin /usr/local
RUN cp -r /protoc/include /usr/local

RUN mkdir /app

WORKDIR /app

# To avoid re-installing dependancies with small changes, perform `go mod download` on build.
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

VOLUME /app
