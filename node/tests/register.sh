#!/bin/sh

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

# The chain has to be running

$CMD/startOneLedger

#status=`$CMD/statusOneLedger`
#
#echo "OneLedger: $status"
#
#if [ -z "$status" ]; then
#	echo "OneLedger isn't running"
#fi

list="Admin Alice Bob Carol"

for name in $list 
do
	address=`$CMD/lookup $name RPCAddress tcp://127.0.0.1:`
	nodeName=`$CMD/lookup $name NodeName`
	WORK=$OLDATA/$nodeName

	$CMD/stopNode $name 

	fullnode register --root $WORK/fullnode --identity $name --address $address
	fullnode register --root $WORK/fullnode --identity $name --address $address --chain OneLedger --pubkey 0x01 --privkey 0x01
	fullnode register --root $WORK/fullnode --identity $name --address $address --chain Bitcoin --pubkey 0x01 --privkey 0x01
	fullnode register --root $WORK/fullnode --identity $name --address $address --chain Ethereum --pubkey 0x01 --privkey 0x01

	$CMD/startNode $name register 
done
