#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

$CMD/startOneLedger

echo "=================== Execute Script =================="
oltest execute --root $OLDATA/David-Node --test hello --fee 0.101 --gas 100
oltest execute --root $OLDATA/Alice-Node --test deadloop --fee 0.101 --gas 100

