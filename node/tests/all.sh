#!/bin/bash

TEST=$GOPATH/src/github.com/Oneledger/protocol/node/tests

$TEST/bigsend.sh
$TEST/connections.sh
$TEST/full.sh
$TEST/query.sh
$TEST/register.sh
$TEST/send.sh
$TEST/simple.sh
$TEST/swap.sh
$TEST/testmint.sh
