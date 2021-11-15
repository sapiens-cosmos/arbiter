#! /bin/bash -x
set -e

go install github.com/regen-network/cosmos-proto/protoc-gen-gocosmos

# Get the path of the cosmos-sdk repo from go/pkg/mod
cosmos_sdk_dir=$(go list -f '{{ .Dir }}' -m github.com/cosmos/cosmos-sdk)
proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  protoc \
  -I "proto" \
  -I "$cosmos_sdk_dir/third_party/proto" \
  -I "$cosmos_sdk_dir/proto" \
  --gocosmos_out=plugins=grpc,\
Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:. \
  --grpc-gateway_out=logtostderr=true,allow_colon_final_segments=true:. \
  $(find "${dir}" -maxdepth 1 -name '*.proto')
done

cp -r github.com/osmosis-labs/osmosis/* ./
rm -rf github.com
