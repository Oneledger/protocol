#!/usr/bin/env

echo "=============================="
echo "Compiling binaries..."
echo "=============================="

curdir=$(pwd)

if [ -z "$OLROOT" ]; then
	source "$GOPATH/src/github.com/Oneledger/protocol/node/scripts/setEnvironment"
fi

OLNODE=$OLROOT/protocol/node

cd $OLNODE

echo "Installing fullnode & olclient..."
make tools && make update && make install

echo "Installing tendermint"
$OLNODE/setup/install_consensus.sh

