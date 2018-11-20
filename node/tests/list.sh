#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

$CMD/startOneLedger

# List out the account details
echo "========== Fullnode Accounts ==========="
olclient list -c "David"
olclient list -c "Alice"
olclient list -c "Bob"
olclient list -c "Carol"
