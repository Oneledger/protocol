#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

$CMD/resetOneLedger
$CMD/startOneLedger

addrAlice=`$CMD/lookup Alice RPCAddress tcp://127.0.0.1:`
addrBob=`$CMD/lookup Bob RPCAddress tcp://127.0.0.1:`

# assumes fullnode is in the PATH
olclient send --party Bob --counterparty Alice --address $addrBob

sleep 3

olclient account -a $addrAlice
olclient account -a $addrBob

$CMD/stopOneLedger
