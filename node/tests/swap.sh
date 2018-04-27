#!/bin/sh

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
OLSCRIPT=$GOPATH/src/github.com/Oneledger/prototype/node/scripts
OLTEST=$GOPATH/src/github.com/Oneledger/prototype/node/tests

# Clear out the existing chains
$OLSCRIPT/resetChain

# Add in the new users
$OLTEST/register.sh

# Startup the chains
$OLSCRIPT/startNode

sleep 9

# Put some money in the user accounts
olclient send -s 1002 --user Admin --to Alice --amount 100000 --currency OLT 
olclient send -s 1003 --user Admin --to Bob --amount 100000 --currency OLT 

# assumes fullnode is in the PATH
olclient swap -s 2001 --user Alice --to 0x0100101010 --amount 3 --currency BTC --exchange 100 --excurrency ETH 
olclient swap -s 2001 --user Bob --to 0x0100101010 --amount 100 --currency ETH --exchange 3 --excurrency BTC 

# Check the balances
olclient account --user Alice
olclient account --user Bob

sleep 3

$OLSCRIPT/stopnode
