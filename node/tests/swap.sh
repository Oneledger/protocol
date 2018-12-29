#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts
TEST=$GOPATH/src/github.com/Oneledger/protocol/node/tests

$TEST/register.sh

echo "================== Test Swap between BTC & ETH ==================="
$CMD/showBalance Alice
sleep 1
$CMD/showBalance Bob

# Let the money get processed
sleep 3

echo "Alice initiate the swap, 5BTC to exchange 100ETH"
olclient swap --root $OLDATA/Alice-Node \
	--party Alice --counterparty Bob --nonce 28 \
	--amount 5.2 --currency BTC --exchange 100.1 --excurrency ETH \
	--fee 0.02 --gas 100

sleep 3

echo "Bob participate the swap 100ETH to exchange 5BTC"
olclient swap --root $OLDATA/Bob-Node  \
	--party Bob --counterparty Alice --nonce 28 \
	--amount 100.1 --currency ETH --exchange 5.2 --excurrency BTC \
	--fee 0.02 --gas 100

echo "Wait for chain to finish"
olclient wait --root $OLDATA/Emma-Node --completed swap --party Alice --party Bob

sleep 5
echo "============================================================="
$CMD/showBalance Alice
sleep 1
$CMD/showBalance Bob
