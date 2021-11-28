#! /bin/bash -x
set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

# Build the docker image first.
docker build --tag arbiter/protogen -f "$DIR"/Dockerfile "$DIR"/../..

# Make sure that previous container not exist.
docker rm --force arbiter_protogen 1>/dev/null 2>&1

docker run -v "$DIR"/../..:/app --name arbiter_protogen arbiter/protogen sh dev/proto/protogen.sh --skip-install

# Clear the container
docker rm --force arbiter_protogen 1>/dev/null 2>&1
