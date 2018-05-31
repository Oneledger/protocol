#!/bin/sh

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
OLTEST=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

$OLTEST/resetChain
$OLTEST/startChain

addrAdmin=`$OLSCRIPT/lookup Admin RPCAddress tcp://127.0.0.1:`
addrAlice=`$OLSCRIPT/lookup Alice RPCAddress tcp://127.0.0.1:`
addrBob=`$OLSCRIPT/lookup Bob RPCAddress tcp://127.0.0.1:`

# assumes fullnode is in the PATH
olclient send --party Bob --counterparty Alice --address $addrBob

sleep 3

$OLTEST/stopChain
