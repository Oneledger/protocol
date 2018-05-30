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

addrAdmin=`$OLSCRIPT/lookup Admin RPCAddress tcp://127.0.0.1:`
addrAlice=`$OLSCRIPT/lookup Alice RPCAddress tcp://127.0.0.1:`
addrBob=`$OLSCRIPT/lookup Bob RPCAddress tcp://127.0.0.1:`

# olclient wait --initialized
#sleep 2 

# Put some money in the user accounts
olclient send --address $addrAdmin -s 1002 --party Admin --counterparty Alice --amount 100000 --currency OLT 
olclient send --address $addrAdmin -s 2003 --party Admin --counterparty Bob --amount 100000 --currency OLT 

# assumes fullnode is in the PATH
olclient swap -s 4001 \
	--address $addrAlice \
	--party Alice --counterparty Bob --nonce 28 \
	--amount 3 --currency BTC --exchange 100 --excurrency ETH 

olclient swap -s 5001 \
	--address $addrBob \
	--party Bob --counterparty Alice --nonce 28 \
	--amount 100 --currency ETH --exchange 3 --excurrency BTC 

olclient wait --completed swap --party Alice --party Bob -s 7001

# Check the balances
olclient account --identity Alice --address $addrAlice 
olclient account --identity Bob --address $addrBob

sleep 3

$OLSCRIPT/stopnode
