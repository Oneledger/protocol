#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

$CMD/startOneLedger

addrAlice=`$OLSCRIPT/lookup Alice RPCAddress tcp://127.0.0.1:`
addrBob=`$OLSCRIPT/lookup Bob RPCAddress tcp://127.0.0.1:`
addrCarol=`$OLSCRIPT/lookup Carol RPCAddress tcp://127.0.0.1:`
addrDavid=`$OLSCRIPT/lookup David RPCAddress tcp://127.0.0.1:`

# Put some money in the user accounts
SEQ=`$CMD/nextSeq`
olclient testmint -s $SEQ -a $addrAlice --party Alice --amount 10000 --currency OLT
olclient testmint -s $SEQ -a $addrBob --party Bob --amount 20000 --currency OLT
olclient testmint -s $SEQ -a $addrCarol --party Carol --amount 3000000 --currency OLT
olclient testmint -s $SEQ -a $addrDavid --party David --amount 1800000 --currency OLT

# assumes fullnode is in the PATH
olclient send -s $SEQ -a $addrBob --party Bob --counterparty Alice --amount 5000 --currency OLT
olclient send -s $SEQ -a $addrAlice --party Alice --counterparty David --amount 5001 --currency OLT
olclient send -s $SEQ -a $addrBob --party Bob --counterparty David --amount 523 --currency OLT
olclient send -s $SEQ -a $addrBob --party Bob --counterparty David --amount 5000 --currency OLT
olclient send -s $SEQ -a $addrAlice --party Alice --counterparty Carol --amount 5000 --currency OLT
olclient send -s $SEQ -a $addrAlice --party Alice --counterparty David --amount 5099 --currency OLT
olclient send -s $SEQ -a $addrDavid --party David --counterparty Carol --amount 51 --currency OLT
olclient send -s $SEQ -a $addrDavid --party David --counterparty Bob --amount 230 --currency OLT
olclient send -s $SEQ -a $addrAlice --party Alice --counterparty Bob --amount 5050 --currency OLT
olclient send -s $SEQ -a $addrBob --party Bob --counterparty Alice --amount 5020 --currency OLT
olclient send -s $SEQ -a $addrAlice --party Alice --counterparty Bob --amount 5050 --currency OLT
olclient send -s $SEQ -a $addrDavid --party David --counterparty Alice --amount 5000 --currency OLT
olclient send -s $SEQ -a $addrAlice --party Alice --counterparty Carol --amount 5099 --currency OLT
olclient send -s $SEQ -a $addrBob --party Bob --counterparty Carol --amount 5012 --currency OLT
olclient send -s $SEQ -a $addrDavid --party David --counterparty Bob --amount 5000 --currency OLT
olclient send -s $SEQ -a $addrDavid --party David --counterparty Alice --amount 45000 --currency OLT
olclient send -s $SEQ -a $addrBob --party Bob --counterparty Carol --amount 543 --currency OLT
olclient send -s $SEQ -a $addrDavid --party David --counterparty Carol --amount 5001 --currency OLT

sleep 8

olclient utxo -a $addrAlice
olclient utxo -a $addrBob
olclient utxo -a $addrCarol
olclient utxo -a $addrDavid

olclient identity -a $addrAlice
olclient identity -a $addrBob
olclient identity -a $addrCarol
olclient identity -a $addrDavid

olclient account -a $addrAlice
olclient account -a $addrBob
olclient account -a $addrCarol
olclient account -a $addrDavid

