#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts
TEST=$GOPATH/src/github.com/Oneledger/protocol/node/tests

# Clear out the existing chains
$CMD/stopOneLedger
$CMD/resetOneLedger

# Add in or update users
$TEST/register.sh
$TEST/send.sh
$TEST/list.sh

$CMD/stopOneLedger

