#!/bin/bash

TEST=$GOPATH/src/github.com/Oneledger/protocol/node/tests

$TEST/accounts.sh
$TEST/bigsend.sh
$TEST/connections.sh
$TEST/full.sh
$TEST/identity.sh
$TEST/query.sh
$TEST/register.sh
$TEST/send.sh
$TEST/simple.sh
$TEST/status.sh
$TEST/swap.sh
$TEST/testmint.sh
$TEST/utxo.sh
