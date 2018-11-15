#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

$CMD/startOneLedger

# Find the addresses
#addrDavid=`$CMD/lookup David RPCAddress tcp://127.0.0.1:`
#addrAlice=`$CMD/lookup Alice RPCAddress tcp://127.0.0.1:`
#addrBob=`$CMD/lookup Bob RPCAddress tcp://127.0.0.1:`
#ddrCarol=`$CMD/lookup Carol RPCAddress tcp://127.0.0.1:`

# List out the account details
echo "========== Fullnode Accounts ==========="
#olclient account -a $addrDavid
olclient list -c "david"
sleep 1
#olclient account -a $addrAlice
olclient list -c "alice"
sleep 1
#olclient account -a $addrBob
olclient list -c "bob"
sleep 1
#olclient account -a $addrCarol
olclient list -c "carol"
sleep 1

