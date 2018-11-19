#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

$CMD/startOneLedger

# List out the account details
echo "========== Fullnode Accounts ==========="
olclient list -c "david"
olclient list -c "alice"
olclient list -c "bob"
olclient list -c "carol"
