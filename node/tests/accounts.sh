#!/bin/sh

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts
OLTEST=$GOPATH/src/github.com/Oneledger/protocol/node/tests

$CMD/startChain

# Find the addresses
addrAdmin=`$CMD/lookup Admin RPCAddress tcp://127.0.0.1:`
addrAlice=`$CMD/lookup Alice RPCAddress tcp://127.0.0.1:`
addrBob=`$CMD/lookup Bob RPCAddress tcp://127.0.0.1:`
addrCarol=`$CMD/lookup Carol RPCAddress tcp://127.0.0.1:`

# List out the account details
olclient account --address $addrAdmin
olclient account --address $addrAlice
olclient account --address $addrBob
olclient account --address $addrCarol
