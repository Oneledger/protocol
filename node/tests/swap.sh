#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts
TEST=$GOPATH/src/github.com/Oneledger/protocol/node/tests

echo "================== Test Swap between BTC & ETH ==================="
$CMD/showBalance Alice
sleep 1
$CMD/showBalance Bob

$TEST/testmint.sh
# Let the money get processed
sleep 3

echo "Alice initiate the swap, 5BTC to exchange 100ETH"
olclient swap -c Alice \
	--party Alice --counterparty Bob --nonce 28 \
	--amount 5 --currency BTC --exchange 100 --excurrency ETH

sleep 3

echo "Bob participate the swap 100ETH to exchange 5BTC"
olclient swap -c Bob \
	--party Bob --counterparty Alice --nonce 28 \
	--amount 100 --currency ETH --exchange 5 --excurrency BTC

echo "Wait for chain to finish"
olclient wait --completed swap --party Alice --party Bob 

sleep 5
echo "============================================================="
$CMD/showBalance Alice
sleep 1
$CMD/showBalance Bob
