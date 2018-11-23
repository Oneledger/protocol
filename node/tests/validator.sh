#!/bin/bash

#
# Test creating a single ApplyValidator transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts
TEST=$GOPATH/src/github.com/Oneledger/protocol/node/tests

echo "================== Test dynamic validator ==================="
$CMD/showBalance Alice
sleep 1
$CMD/showBalance Bob
sleep 1
$CMD/showBalance Emma

$TEST/testmint.sh
# Let the money get processed
sleep 3

echo "Emma sending validator token with 5 coins"
olclient applyvalidator -c Emma --id Emma --amount 5

sleep 3

echo "============================================================="
$CMD/showBalance Alice
sleep 1
$CMD/showBalance Bob
sleep 1
$CMD/showBalance Emma
