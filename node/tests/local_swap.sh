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

addrAlice=`$CMD/lookup Alice RPCAddress tcp://127.0.0.1:`
addrBob=`$CMD/lookup Bob RPCAddress tcp://127.0.0.1:`

# Put some money in the user accounts
SEQ=`$CMD\nextSeq`
olclient send -s $SEQ -a $addrAlice --party Alice --amount 100000 --currency OLT 

SEQ=`$CMD\nextSeq`
olclient send -s $SEQ -a $addrBob --party Bob --amount 100000 --currency OLT 

# assumes fullnode is in the PATH
SEQ=`$CMD\nextSeq`
olclient swap -s $SEQ -a $addrAlice \
	--party Alice --counterparty Bob --nonce 28 \
	--amount 3 --currency OLT --exchange 100 --excurrency ETH 

SEQ=`$CMD\nextSeq`
olclient swap -s $SEQ -a $addrBob \
	--party Bob --counterparty Alice --nonce 28 \
	--amount 100 --currency ETH --exchange 3 --excurrency OLT 

# Check the balances
olclient account -a $addrAlice --identity Alice
olclient account -a $addrBob --identity Bob

sleep 3

$CMD/stopOneLedger
