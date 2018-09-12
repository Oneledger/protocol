#!/usr/bin/env bash

mkdir -p $OLROOT/protocol/VOLUME

for binary in fullnode tendermint olclient
do
	binaryPath=`whereis ${binary} | awk '{print $2}'`

	if [ -z $binaryPath ]
	then
		echo "Couldn't find $binaryPath in PATH"
	else
		cp -v $binaryPath $OLROOT/protocol/VOLUME
	fi
done
