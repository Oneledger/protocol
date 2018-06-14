#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

# Startup the chains
$CMD/startOneLedger

# Find the addresses
addrAdmin=`$CMD/lookup Admin RPCAddress tcp://127.0.0.1:`
addrAlice=`$CMD/lookup Alice RPCAddress tcp://127.0.0.1:`
addrBob=`$CMD/lookup Bob RPCAddress tcp://127.0.0.1:`

# List out the account details
olclient identity -a $addrAlice --identity Alice
olclient account  -a $addrAlice --identity Alice

$CMD/stopOneLedger
