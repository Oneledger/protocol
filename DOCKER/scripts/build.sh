#!/usr/bin/env bash

echo "=============================="
echo "Compiling binaries..."
echo "=============================="

if [ -z "$OLROOT" ]; then
	source "./node/scripts/setEnvironment"
fi

OLNODE=$OLROOT/protocol/node

cd $OLNODE

make tools && make update && make install

$OLNODE/setup/install_consensus.sh
