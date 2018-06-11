#!/bin/sh

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts
TEST=$GOPATH/src/github.com/Oneledger/protocol/node/tests

# Clear out the existing chains
$CMD/resetOneLedger

# Add in or update users
$TEST/register.sh

# Startup the chains
$CMD/startOneLedger

addrAdmin=`$CMD/lookup Admin RPCAddress tcp://127.0.0.1:`
addrAlice=`$CMD/lookup Alice RPCAddress tcp://127.0.0.1:`
addrBob=`$CMD/lookup Bob RPCAddress tcp://127.0.0.1:`

# olclient wait --initialized
#sleep 2 

# Put some money in the user accounts
SEQ=`$CMD/nextSeq`
olclient send $SEQ -a $addrAdmin --party Admin --counterparty Alice --amount 100000 --currency OLT 

SEQ=`$CMD/nextSeq`
olclient send $SEQ -a $addrAdmin --party Admin --counterparty Bob --amount 100000 --currency OLT 

# Let the money get processed
sleep 3

# assumes fullnode is in the PATH
SEQ=`$CMD/nextSeq`
olclient swap $SEQ -a $addrAlice \
	--party Alice --counterparty Bob --nonce 28 \
	--amount 3 --currency BTC --exchange 100 --excurrency ETH 

olclient swap $SEQ -a $addrBob \
	--party Bob --counterparty Alice --nonce 28 \
	--amount 100 --currency ETH --exchange 3 --excurrency BTC 


# Wait for the swaps to complete
olclient wait --completed swap --party Alice --party Bob 

# Check the final balances
olclient account -a $addrAlice --identity Alice
olclient account -a $addrBob --identity Bob

$CMD/stopOneLedger
