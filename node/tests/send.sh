#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

$CMD/startOneLedger

addrAlice=`$OLSCRIPT/lookup Alice RPCAddress tcp://127.0.0.1:`
addrBob=`$OLSCRIPT/lookup Bob RPCAddress tcp://127.0.0.1:`
addrEmma=`$OLSCRIPT/lookup Emma RPCAddress tcp://127.0.0.1:`


echo "=================== Test Send Transaction =================="
# Put some money in the user accounts
SEQ=`$CMD/nextSeq`
olclient testmint -s $SEQ -a $addrAlice --party Alice --amount 10000 --currency OLT

sleep 3

olclient testmint -s $SEQ -a $addrBob --party Bob --amount 25000 --currency OLT

sleep 3

olclient testmint -s $SEQ -a $addrEmma --party Emma --amount 5000 --currency OLT

sleep 3

echo "Finished Minting"

sleep 10

# assumes fullnode is in the PATH
olclient send -s $SEQ -a $addrBob --party Bob --counterparty Alice --amount 5000 --currency OLT

sleep 6

olclient account -a $addrBob

olclient account -a $addrAlice
