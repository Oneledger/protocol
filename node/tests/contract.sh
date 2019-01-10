#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node
TESTS=$CMD/olvm/interpreter/samples


echo "=================== Test Contracts =================="
olclient install --root $OLDATA/David-Node --owner David-OneLedger --name Test --version v0.0.1 -f $TESTS/helloworld.js --fee 0.109
#olclient install --root $OLDATA/David-Node --owner David-OneLedger --name Test --version v0.0.1 -f $CMD/tests/data/contract1.txt
#olclient install --root $OLDATA/David-Node --owner David-OneLedger --name Test --version v0.0.2 -f $CMD/tests/data/contract2.txt
#olclient install --root $OLDATA/David-Node --owner David-OneLedger --name Test2 --version v0.0.1 -f $CMD/tests/data/contract3.txt

olclient execute --root $OLDATA/David-Node --owner David-OneLedger --name Test --version v0.0.1 --fee 0.2001 --gas 100
#olclient execute --root $OLDATA/David-Node --owner David-OneLedger --name Test --version v0.0.2
#olclient execute --root $OLDATA/David-Node --owner David-OneLedger --name Test2 --version v0.0.1

sleep 10
$CMD/scripts/stopDevNet

