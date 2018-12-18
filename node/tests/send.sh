#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

$CMD/startOneLedger

echo "=================== Send Transactions =================="
olclient send -c Bob --party Bob --counterparty Alice --amount 5000.112 --currency OLT --fee 1.01
olclient send -c Bob --party Bob --counterparty Carol --amount 2000.234 --currency OLT --fee 5.03

sleep 6
