#! /bin/bash -x
set -e

PROTOC_URL=https://github.com/protocolbuffers/protobuf/releases/download/v3.19.1/protoc-3.19.1-linux-x86_64.zip
PROTOC_FILENAME=protoc-3.19.1-linux-x86_64.zip

if [ "aarch64" = $(uname -m) ]; then
  PROTOC_URL=https://github.com/protocolbuffers/protobuf/releases/download/v3.19.1/protoc-3.19.1-linux-aarch_64.zip
  PROTOC_FILENAME=protoc-3.19.1-linux-aarch_64.zip
fi

wget $PROTOC_URL
unzip $PROTOC_FILENAME -d /protoc
cp -r /protoc/bin /usr/local
cp -r /protoc/include /usr/local
