#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

$CMD/startOneLedger

addrAlice=`$OLSCRIPT/lookup Alice RPCAddress tcp://127.0.0.1:`
addrBob=`$OLSCRIPT/lookup Bob RPCAddress tcp://127.0.0.1:`

echo "===================test send transaction=================="
# Put some money in the user accounts
SEQ=`$CMD/nextSeq`
olclient testmint -s $SEQ -a $addrAlice --party Alice --amount 10000 --currency OLT
olclient testmint -s $SEQ -a $addrBob --party Bob --amount 20000 --currency OLT

# assumes fullnode is in the PATH
olclient send -s $SEQ -a $addrBob --party Bob --counterparty Alice --amount 5000 --currency OLT

sleep 10

olclient account -a $addrBob

sleep 1

olclient account -a $addrAlice

sleep 3
