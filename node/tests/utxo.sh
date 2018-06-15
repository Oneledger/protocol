#!/bin/sh

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

$CMD/startOneLedger

# Find the addresses
addrDavid=`$CMD/lookup David RPCAddress tcp://127.0.0.1:`
addrAlice=`$CMD/lookup Alice RPCAddress tcp://127.0.0.1:`
addrBob=`$CMD/lookup Bob RPCAddress tcp://127.0.0.1:`
addrCarol=`$CMD/lookup Carol RPCAddress tcp://127.0.0.1:`

# List out the account details
olclient utxo -a $addrDavid
olclient utxo -a $addrAlice
olclient utxo -a $addrBob
olclient utxo -a $addrCarol
