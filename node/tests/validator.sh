#!/bin/bash

#
# Test creating a single ApplyValidator transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts
TEST=$GOPATH/src/github.com/Oneledger/protocol/node/tests


$TEST/register.sh

echo "================== Test dynamic validator ==================="

olclient testmint --root $OLDATA/Emma-Node --party Emma --amount 10 --currency VT
# Let the money get processed
sleep 3

echo "Emma sending validator token with 5 coins"
olclient applyvalidator --root $OLDATA/Emma-Node  --id Emma --amount 5

sleep 4

olclient testmint --root $OLDATA/David-Node --party David --amount 10 --currency VT

sleep 4

olclient send --root $OLDATA/David-Node --party David --counterparty Bob --amount 6 --currency VT

echo "============================================================="
olclient list --root $OLDATA/Emma-Node
olclient list --root $OLDATA/David-Node

num=`pgrep -f "^olfullnode node" | wc -l `

if [[ $num < 5 ]]; then
    echo "Validator test failed"
    exit  -1
fi

echo "Validator test succeed"

$CMD/stopOneLedger
