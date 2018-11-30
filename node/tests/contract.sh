#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node


echo "=================== Test Contracts =================="
olclient install -c David --owner David-OneLedger --name Test --version v0.0.1 -f $CMD/tests/data/contract1.txt
#olclient install -c David --owner David-OneLedger --name Test --version v0.0.2 -f $CMD/tests/data/contract2.txt
#olclient install -c David --owner David-OneLedger --name Test2 --version v0.0.1 -f $CMD/tests/data/contract3.txt

olclient execute -c David --owner David-OneLedger --name Test --version v0.0.1
#olclient execute -c David --owner David-OneLedger --name Test --version v0.0.2
#olclient execute -c David --owner David-OneLedger --name Test2 --version v0.0.1

sleep 10
$CMD/scripts/stopDevNet

