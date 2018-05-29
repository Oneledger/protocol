#!/bin/sh

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
OLSCRIPT=$GOPATH/src/github.com/Oneledger/protocol/node/scripts
OLTEST=$GOPATH/src/github.com/Oneledger/protocol/node/tests

# Clear out the existing chains
$OLSCRIPT/resetChain

# Add in or update users
$OLTEST/register.sh

# Startup the chains
$OLSCRIPT/startChain

# olclient wait --initialized
sleep 6 

# Put some money in the user accounts
olclient send -s 1002 --party Admin --counterparty Alice --amount 100000 --currency OLT 
olclient send -s 1003 --party Admin --counterparty Bob --amount 100000 --currency OLT 

# assumes fullnode is in the PATH
olclient swap -s 2001 \
	--address tcp://127.0.0.1:46612 \
	--party Alice --counterparty Bob --nonce 28 \
	--amount 3 --currency OLT --exchange 100 --excurrency ETH 

olclient swap -s 2001 \
	--address tcp://127.0.0.1:46613 \
	--party Bob --counterparty Alice --nonce 28 \
	--amount 100 --currency ETH --exchange 3 --excurrency OLT 

olclient wait --completed swap --party Alice --party Bob -s 2001

# Check the balances
olclient account --user Alice
olclient account --user Bob

sleep 3

$OLSCRIPT/stopnode
