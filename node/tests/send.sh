#!/bin/sh

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
OLTEST=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

$OLTEST/startOneLedger

addrAdmin=`$OLSCRIPT/lookup Admin RPCAddress tcp://127.0.0.1:`
addrAlice=`$OLSCRIPT/lookup Alice RPCAddress tcp://127.0.0.1:`
addrBob=`$OLSCRIPT/lookup Bob RPCAddress tcp://127.0.0.1:`

# Put some money in the user accounts
olclient send --address $addrAdmin -s 1002 --party Admin --counterparty Alice --amount 10000 --currency OLT 
olclient send --address $addrAdmin -s 2003 --party Admin --counterparty Bob --amount 20000 --currency OLT 

# assumes fullnode is in the PATH
olclient send --address $addrBob -s 3004 --party Bob --counterparty Alice --amount 5000 --currency OLT

sleep 3

$OLTEST/stopOneLedger
