#!/usr/bin/env sh

ROOT_PATH=$(cd $(dirname $BASH_SOURCE[0])/.. && pwd)
cd $ROOT_PATH

docker build --no-cache -t dhy1210/bitcoin-core -f Dockerfile.bitcoind-core .

