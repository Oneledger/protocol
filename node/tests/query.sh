#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

# Startup the chains
$CMD/startOneLedger

olclient list -c Alice 
olclient list -c Bob 

$CMD/stopOneLedger
