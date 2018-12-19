#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

# Startup the chains
$CMD/startOneLedger

olclient list --root $OLDATA/Bob-Node
olclient list --root $OLDATA/Alice-Node

$CMD/stopOneLedger
