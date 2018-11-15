#!/bin/bash
#
# Register all of the identities and accounts on OneLedger
#
# Need to test to see if this has already been done...
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

list="David Alice Bob Carol Emma"
#list="David"

$CMD/startOneLedger

echo "=================== Test Registration ======================="
for name in $list
do
	nodeAddr=`$CMD/lookup $name RPCAddress tcp://127.0.0.1:`
	nodeName=`$CMD/lookup $name NodeName`
	WORK=$OLDATA/$nodeName
	DATA=$WORK/tendermint
	LOG=$WORK
	ROOT=$WORK/olfullnode

	echo "Register [$name] "
	$CMD/stopNode $name

	# Setup a global Identity and OneLedger account
	olfullnode register --root $ROOT -a $nodeAddr \
		--node $nodeName \
		--identity $name \
		--tendermintRoot $DATA \
		>> $LOG/olfullnode.log 2>&1

	# Fill in the specific chain accounts
	olfullnode register --root $ROOT -a $nodeAddr \
		--node $nodeName \
		--identity $name --chain Bitcoin \
		--tendermintRoot $DATA \
		>> $LOG/olfullnode.log 2>&1

	olfullnode register --root $ROOT -a $nodeAddr \
		--node $nodeName \
		--identity $name --chain Ethereum \
		--tendermintRoot $DATA \
		>> $LOG/olfullnode.log 2>&1

	# Broadtcast it to all of the nodes to make sure it is unique
	$CMD/startNode $name register

	# Need to let the identity transaction fully broadcast, before letting the next node shutdown.
	sleep 10
done
