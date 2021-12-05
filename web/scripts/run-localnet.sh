#! /bin/bash -x
set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

echo "$DIR"/../localnet/Dockerfile "$DIR"/../..

# Build the docker image first.
docker build --tag arbiter/localnet -f "$DIR"/../localnet/Dockerfile "$DIR"/../..

# Make sure that previous container not exist.
docker rm --force arbiter_localnet

# Start container as daemon with some ports opening.
docker run -d -p 1317:1317 -p 26657:26657 -p 9090:9090 --name arbiter_localnet arbiter/localnet

echo "Validator mnemonic: high gain deposit chuckle hundred regular exist approve peanut enjoy comfort ride"
echo "Account1 mnemonic: health nest provide snow total tissue intact loyal cargo must credit wrist"
echo "Account2 mnemonic: canyon stone next tenant trial ugly slim luggage ski govern outside comfort"
echo "Account2 mnemonic: travel renew first fiction trick fly disease advance hunt famous absurd region"
echo "Each account has the balances (10000000000uarb,10000000000ugreen)"

echo "Chain id: arbiter-localnet-1"
echo "RPC: http://localhost:26657"
echo "LCD: http://localhost:1317"
echo "Swagger doc: http://localhost:1317/swagger/#"

echo "Docker container is running on \"arbiter_localnet\""
echo "After testing, to remove existing container, run \"sudo docker rm --force arbiter_localnet\""
