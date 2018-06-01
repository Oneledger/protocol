#!/bin/sh

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
OLTEST=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

status=`$OLTEST/statusChain`

echo "ChainStatus: $status"

if [ -z "$status" ]; then
	echo "Chain isn't running"
fi

list="Admin Alice Bob Carol"

for name in $list 
do
	address=`$OLSCRIPT/lookup $name RPCAddress tcp://127.0.0.1:`

	$OLTEST/stopNode $name 

	fullnode register --identity $name --address $address
	fullnode register --identity $name --address $address --chain OneLedger --pubkey 0x01 --privkey 0x01
	fullnode register --identity $name --address $address --chain Bitcoin --pubkey 0x01 --privkey 0x01
	fullnode register --identity $name --address $address --chain Ethereum --pubkey 0x01 --privkey 0x01

	$OLTEST/startNode $name register 
done
