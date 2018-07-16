#!/bin/bash

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

# olclient wait --initialized
#sleep 2 
echo "==================Test swap between BTC & ETH==================="


# Put some money in the user accounts
SEQ=`$CMD/nextSeq`
olclient testmint -s $SEQ -a $addrAlice --party Alice --amount 100000 --currency OLT 

SEQ=`$CMD/nextSeq`
olclient testmint -s $SEQ -a $addrBob --party Bob --amount 100000 --currency OLT 


# Let the money get processed
sleep 3

echo "Alice initiate the swap, 5BTC to exchange 100ETH"
# assumes fullnode is in the PATH
SEQ=`$CMD/nextSeq`
olclient swap -s $SEQ -a $addrAlice \
	--party Alice --counterparty Bob --nonce 28 \
	--amount 5 --currency BTC --exchange 100 --excurrency ETH

sleep 3

echo "Bob participate the swap 100ETH to exchange 5BTC"
olclient swap -s $SEQ -a $addrBob \
	--party Bob --counterparty Alice --nonce 28 \
	--amount 100 --currency ETH --exchange 5 --excurrency BTC

echo "Wait for chain to finish"
# Wait for the swaps to complete
olclient wait --completed swap --party Alice --party Bob 

sleep 5
echo "============================================================="
#echo "Check the account balance"
## Check the final balances
#olclient account -a $addrAlice
#
#sleep 3
#
#olclient account -a $addrBob
#
#sleep 1

$CMD/stopOneLedger
