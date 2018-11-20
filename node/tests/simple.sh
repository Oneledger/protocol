#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

$CMD/resetOneLedger
$CMD/startOneLedger

olclient send -c Bob --party Bob --counterparty Alice --amount 1000 --currency OLT

sleep 3

$CMD/stopOneLedger
