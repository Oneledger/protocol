#!/bin/bash

#
# Test creating a single ApplyValidator transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts
TEST=$GOPATH/src/github.com/Oneledger/protocol/node/tests

echo "================== Test dynamic validator ==================="
$TEST/register.sh
$TEST/testmint.sh

olclient testmint -c Emma --party Emma --amount 10 --currency VT
# Let the money get processed
sleep 3

echo "Emma sending validator token with 5 coins"
olclient applyvalidator -c Emma --id Emma --amount 5

sleep 6

echo "============================================================="
olclient list -c Emma


$CMD/stopOneLedger
