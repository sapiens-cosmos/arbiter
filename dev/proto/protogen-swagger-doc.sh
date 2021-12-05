#! /bin/bash -x
set -e

mkdir -p ./tmp-swagger-gen

# Get the path of the cosmos-sdk repo from go/pkg/mod
cosmos_sdk_dir=$(go list -f '{{ .Dir }}' -m github.com/cosmos/cosmos-sdk)
proto_dirs=$(find ./proto "$cosmos_sdk_dir"/proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do

  # generate swagger files (filter query files)
  query_file=$(find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \))

  if [ -n "$query_file" ]; then
    protoc \
      -I "proto" \
      -I "$cosmos_sdk_dir/third_party/proto" \
      -I "$cosmos_sdk_dir/proto" \
      "$query_file" \
      --swagger_out=./tmp-swagger-gen \
      --swagger_opt=logtostderr=true --swagger_opt=fqn_for_swagger_name=true --swagger_opt=simple_operation_ids=true
  fi
done

rm ./client/docs/swagger_legacy.yaml
cp "$cosmos_sdk_dir"/client/docs/swagger_legacy.yaml ./client/docs/swagger_legacy.yaml

# combine swagger files
# uses nodejs package `swagger-combine`.
# all the individual swagger files need to be configured in `config.json` for merging
swagger-combine ./client/docs/config.json -o ./client/docs/swagger-ui/swagger.yaml -f yaml --continueOnConflictingPaths true --includeDefinitions true

# clean swagger files
rm -rf ./tmp-swagger-gen

statik -src=./client/docs/swagger-ui -dest=./client/docs -f -m
