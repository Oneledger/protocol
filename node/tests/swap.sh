#!/bin/sh

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
OLSCRIPT=$GOPATH/src/github.com/Oneledger/prototype/node/scripts
OLTEST=$GOPATH/src/github.com/Oneledger/prototype/node/tests

# Clear out the existing chains
$OLSCRIPT/resetChain

# Add in the new users
$OLTEST/register.sh

# Startup the chains
$OLSCRIPT/startNode

sleep 9

# assumes fullnode is in the PATH
fullnode swap --user Alice --from 0x01010100101 --to 0x0100101010 --amount 3BTC --amount 100ETH 
fullnode swap --user Bob --from 0x01010100101 --to 0x0100101010 --amount 3BTC --amount 100ETH 

fullnode account --user Alice

sleep 3

$OLSCRIPT/stopnode
