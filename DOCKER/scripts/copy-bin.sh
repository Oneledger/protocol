#!/usr/bin/env bash

if [ -z "$OLROOT" ]; then
	OLROOT=$(pwd)
fi

mkdir -p $OLROOT/protocol/VOLUME

for binary in fullnode tendermint olclient
do
	binaryPath=`whereis ${binary} | awk '{print $2}'`

	if [ -z $binaryPath ]
	then
		echo "Couldn't find $binary in PATH. Try running \`make build\` first."
		exit 1
	else
		cp $binaryPath $OLROOT/protocol/VOLUME
	fi
done
