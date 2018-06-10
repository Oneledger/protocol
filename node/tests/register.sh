#!/bin/sh

#
# Register all of the identities and accounts on OneLedger
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

list="Admin Alice Bob Carol"

for name in $list 
do
	address=`$CMD/lookup $name RPCAddress tcp://127.0.0.1:`
	nodeName=`$CMD/lookup $name NodeName`
	WORK=$OLDATA/$nodeName

	$CMD/stopNode $name 

	# Setup an Identity
	fullnode register --root $WORK/fullnode --identity $name --address $address

	# Associated it with a OneLedger account
	fullnode register --root $WORK/fullnode --identity $name --address $address \
		--chain OneLedger --pubkey 0x01 --privkey 0x01

	# Broadtcast it to all of the nodes to make sure it is unique
	$CMD/startNode $name register 
	sleep 5
	$CMD/stopNode $name 

	# Fill in the specific chain accounts
	fullnode register --root $WORK/fullnode --identity $name --address $address \
		--chain Bitcoin --pubkey 0x01 --privkey 0x01

	fullnode register --root $WORK/fullnode --identity $name --address $address \
		--chain Ethereum --pubkey 0x01 --privkey 0x01

	# Everything should be functional now
	$CMD/startNode $name 
done
