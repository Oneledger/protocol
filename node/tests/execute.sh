#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

$CMD/startOneLedger

echo "=================== Execute Script =================="
oltest execute -c David --test hello --fee 0.01 --gas 0.10
oltest execute -c Alice --test deadloop --fee 0.01 --gas 0.10

