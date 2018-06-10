#!/bin/sh

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

$CMD/startOneLedger

addrAdmin=`$OLSCRIPT/lookup Admin RPCAddress tcp://127.0.0.1:`
addrAlice=`$OLSCRIPT/lookup Alice RPCAddress tcp://127.0.0.1:`
addrBob=`$OLSCRIPT/lookup Bob RPCAddress tcp://127.0.0.1:`

# Put some money in the user accounts
olclient send -a $addrAdmin -s 1002 --party Admin --counterparty Alice --amount 10000 --currency OLT 
olclient send -a $addrAdmin -s 2003 --party Admin --counterparty Bob --amount 20000 --currency OLT 

# assumes fullnode is in the PATH
olclient send -a $addrBob -s 3004 --party Bob --counterparty Alice --amount 5000 --currency OLT

sleep 5

olclient account -a $addrAdmin
olclient account -a $addrBob
olclient account -a $addrAlice

