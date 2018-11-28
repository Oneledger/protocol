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
	# Add the accounts, keys are generated internally
	olclient update -c $name --account "$name-OneLedger"

	# todo: need to flag to set node account for the node.
#	olclient update -c $name --account "$name-BitCoin" --chain "BitCoin"
#	olclient update -c $name --account "$name-Ethereum" --chain "Ethereum"

	olclient register -c $name --identity "$name" --account "$name-OneLedger" --node "$name-Node"
done

# Give it some time to get committed
#sleep 20
olclient wait --completed identity --identity $list
