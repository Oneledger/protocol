#!/bin/sh

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
OLTEST=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

status=`$OLTEST/statusChain`

echo "ChainStatus: $status"

if [ -z "$status" ]; then
	echo "Chain wasn't running"
else
	$OLTEST/stopChain
fi

list="Admin Alice Bob Carol"

for name in $list 
do
	address=`$OLSCRIPT/lookup $name RPCAddress tcp://127.0.0.1:`

	fullnode register --address $address --identity $name

	fullnode register --address $address --chain OneLedger --identity $name \
		--pubkey 0x0103a39e93332 --privkey 0x0103a39e93332

	fullnode register --address $address --chain Bitcoin --identity $name \
		--pubkey 0x0203a39e93332 --privkey 0x0203a39e93332 

	fullnode register --address $address --chain Ethereum --identity $name \
		--pubkey 0x0303a39e93332 --privkey 0x0303a39e93332 
done

if [ -z "$status" ]; then
	echo "Chain isn't restarted"
else
	$OLTEST/startChain
fi
