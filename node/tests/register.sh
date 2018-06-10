#!/bin/sh
#
# Register all of the identities and accounts on OneLedger
#
# Need to test to see if this has already been done...
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

list="Admin Alice Bob Carol"

for name in $list 
do
	nodeAddr=`$CMD/lookup $name RPCAddress tcp://127.0.0.1:`
	nodeName=`$CMD/lookup $name NodeName`
	WORK=$OLDATA/$nodeName
	ROOT=$WORK/fullnode

	$CMD/stopNode $name 

	# Setup a global Identity
	fullnode register --root $ROOT -a $nodeAddr \
		--identity $name 

	# Associated it with a OneLedger account
	fullnode register --root $ROOT -a $nodeAddr \
		--identity $name --chain OneLedger --pubkey 0x01010100111 --privkey 0x01

	# Broadtcast it to all of the nodes to make sure it is unique
	$CMD/startNode $name register 
	sleep 5
	$CMD/stopNode $name 

	# Fill in the specific chain accounts
	fullnode register --root $ROOT -a $nodeAddr \
		--identity $name --chain Bitcoin --pubkey 0x01 --privkey 0x01

	fullnode register --root $ROOT -a $nodeAddr \
		--identity $name --chain Ethereum --pubkey 0x01 --privkey 0x01

	# Everything should be functional now
	$CMD/startNode $name 
done
