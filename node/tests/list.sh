#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

$CMD/startOneLedger

# List out the account details
echo "========== Fullnode Accounts ==========="
olclient list --root $OLDATA/David-Node
olclient list --root $OLDATA/Alice-Node
olclient list --root $OLDATA/Bob-Node
olclient list --root $OLDATA/Carol-Node
olclient list --root $OLDATA/Emma-Node
