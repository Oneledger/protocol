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

olclient account --user Alice

sleep 3

$OLSCRIPT/stopnode
