#!/bin/sh

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
OLSCRIPT=$GOPATH/src/github.com/Oneledger/protocol/node/scripts
OLTEST=$GOPATH/src/github.com/Oneledger/protocol/node/tests

# Clear out the existing chains
$OLSCRIPT/resetChain

# Add in the new users
$OLTEST/register.sh

# Startup the chains
$OLSCRIPT/startChain

olclient account --identity Alice --address tcp://127.0.0.1:46601

$OLSCRIPT/stopnode
