#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

$CMD/startOneLedger


echo "=================== Send Transactions =================="
olclient send -c Bob --party Bob --counterparty Alice --amount 5000 --currency OLT
olclient send -c Bob --party Bob --counterparty Carol --amount 2000 --currency OLT

sleep 6
