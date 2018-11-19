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
addrCarol=`$CMD/lookup Carol RPCAddress tcp://127.0.0.1:`
addrDavid=`$CMD/lookup David RPCAddress tcp://127.0.0.1:`
addrEmma=`$CMD/lookup Emma RPCAddress tcp://127.0.0.1:`

# olclient wait --initialized
#sleep 2 

# Put some money in the user accounts
SEQ=`$CMD/nextSeq`
olclient testmint -s $SEQ -a $addrAlice --party Alice --amount 100001 --currency OLT 

SEQ=`$CMD/nextSeq`
olclient testmint -s $SEQ -a $addrBob --party Bob --amount 50002 --currency OLT 

SEQ=`$CMD/nextSeq`
olclient testmint -s $SEQ -a $addrCarol --party Carol --amount 25003 --currency OLT 

SEQ=`$CMD/nextSeq`
olclient testmint -s $SEQ -a $addrDavid --party David --amount 12004 --currency OLT

SEQ=`$CMD/nextSeq`
olclient testmint -s $SEQ -a $addrEmma --party Emma --amount 50005 --currency OLT

echo "Finished Minting"

sleep 10

olclient account -a $addrAlice 
olclient account -a $addrBob 

$CMD/stopOneLedger
