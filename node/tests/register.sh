#!/bin/sh

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
OLTEST=$GOPATH/src/github.com/Oneledger/prototype/node/scripts

# assumes fullnode is in the PATH
fullnode register --chain OneLedger --name Paul --pubkey 0x01 --privkey 0x01
fullnode register --chain Bitcoin --name Paul --pubkey 0x02 --privkey 0x02
fullnode register --chain Ethereum --name Paul --pubkey 0x02 --privkey 0x02

