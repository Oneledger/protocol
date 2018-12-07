#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts
TEST=$GOPATH/src/github.com/Oneledger/protocol/node/tests

# Startup the chains
#$CMD/startOneLedger

echo "=================== Testmint Transactions =================="

# Put some money in the user accounts
olclient testmint -c Alice --party Alice --amount 10000 --currency OLT
olclient testmint -c Bob --party Bob --amount 50000 --currency OLT
olclient testmint -c Carol --party Carol --amount 25000 --currency OLT
olclient testmint -c David --party David --amount 12000 --currency OLT

sleep 10
echo "Finished Minting"
