#!/bin/sh

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
OLTEST=$GOPATH/src/github.com/Oneledger/prototype/node/scripts

# assumes fullnode is in the PATH
#olclient register --address $ADDRESS --name Paul 
#olclient register --address $ADDRESS --chain OneLedger --name Paul --pubkey 0x01 --privkey 0x01
#olclient register --address $ADDRESS --chain Bitcoin --name Paul --pubkey 0x02 --privkey 0x02
#olclient register --address $ADDRESS --chain Ethereum --name Paul --pubkey 0x02 --privkey 0x02

# TODO: Stop the chain, add the accounts, then restart

ADDRESS="tcp://127.0.0.1:46621"

status=`$OLTEST/statusChain`

echo "ChainStatus: $status"

if [ -z "$status" ]; then
	echo "Chain wasn't running"
else
	$OLTEST/stopChain
fi

fullnode register --address $ADDRESS --name Paul
fullnode register --address $ADDRESS --chain OneLedger --name Paul --pubkey 0x01 --privkey 0x01
fullnode register --address $ADDRESS --chain Bitcoin --name Paul --pubkey 0x01 --privkey 0x01
fullnode register --address $ADDRESS --chain Ethereum --name Paul --pubkey 0x01 --privkey 0x01

if [ -z "$status" ]; then
	echo "Chain isn't restarted"
else
	$OLTEST/startChain
fi
