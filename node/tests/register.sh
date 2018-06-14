#!/bin/bash
#
# Register all of the identities and accounts on OneLedger
#
# Need to test to see if this has already been done...
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

list="Alice Bob Carol"

$CMD/startOneLedger

for name in $list 
do
	nodeAddr=`$CMD/lookup $name RPCAddress tcp://127.0.0.1:`
	nodeName=`$CMD/lookup $name NodeName`
	WORK=$OLDATA/$nodeName
	ROOT=$WORK/fullnode

	$CMD/stopNode $name 

	# Setup a global Identity and OneLedger account
	fullnode register --root $ROOT -a $nodeAddr \
		--node $nodeName \
		--identity $name 

	# Fill in the specific chain accounts
	fullnode register --root $ROOT -a $nodeAddr \
		--node $nodeName \
		--identity $name --chain Bitcoin 

	fullnode register --root $ROOT -a $nodeAddr \
		--node $nodeName \
		--identity $name --chain Ethereum 

	# Broadtcast it to all of the nodes to make sure it is unique
	$CMD/startNode $name register 
	sleep 5
done
