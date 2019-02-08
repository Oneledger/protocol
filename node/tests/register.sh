#!/bin/bash
#
# Register all of the identities and accounts on OneLedger
#
# Need to test to see if this has already been done...
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

list="David Alice Bob Carol Emma"

$CMD/startOneLedger

echo "=================== Test Registration ======================="
for name in $list
do
	# Add the accounts, keys are generated internally
	olclient update --root $OLDATA/$name-Node --account "$name-OneLedger"

	#TODO: Need flag to set node account for the each external chain.
	olclient update --root $OLDATA/$name-Node --account "$name-BitCoin" --chain "BitCoin"
	olclient update --root $OLDATA/$name-Node --account "$name-Ethereum" --chain "Ethereum"

	# Account must have money in order to pay the registration fee
	olclient testmint --root $OLDATA/$name-Node --party "$name-OneLedger" --amount 100000.1 --currency OLT

	# Register the identity across the chain
	olclient register --root $OLDATA/$name-Node --identity "$name" \
		--account "$name-OneLedger" --node "$name-Node" --fee 0.1
done

# Give it some time to get committed
olclient --root $OLDATA/Emma-Node wait --completed identity --identity $list
