#!/bin/sh

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
OLTEST=$GOPATH/src/github.com/Oneledger/prototype/node/scripts

$OLTEST/resetChain

$OLTEST/startNode

sleep 9

# assumes fullnode is in the PATH
fullnode send --to 0x01010100101

sleep 3

$OLTEST/stopnode
