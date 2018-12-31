#!/bin/bash

# Stop, reset, start and register, so a test can be performed.
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

$CMD/stopOneLedger
$CMD/resetOneLedger
$CMD/startOneLedger

TESTS=$GOPATH/src/github.com/Oneledger/protocol/node/tests
$TESTS/register.sh

