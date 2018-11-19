#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

$CMD/startOneLedger

# assumes olfullnode is in the PATH
SEQ=`$CMD/nextSeq`
olclient send -c Bob --party Bob --counterparty Alice --amount 5000 --currency OLT
olclient send -c Alice --party Alice --counterparty David --amount 5001 --currency OLT
olclient send -c Bob --party Bob --counterparty David --amount 523 --currency OLT
olclient send -c Bob --party Bob --counterparty David --amount 5000 --currency OLT
olclient send -c Alice --party Alice --counterparty Carol --amount 5000 --currency OLT
olclient send -c Alice --party Alice --counterparty David --amount 5099 --currency OLT
olclient send -c David --party David --counterparty Carol --amount 51 --currency OLT
olclient send -c David --party David --counterparty Bob --amount 230 --currency OLT
olclient send -c Alice --party Alice --counterparty Bob --amount 5050 --currency OLT
olclient send -c Bob --party Bob --counterparty Alice --amount 5020 --currency OLT
olclient send -c Alice --party Alice --counterparty Bob --amount 5050 --currency OLT
olclient send -c David --party David --counterparty Alice --amount 5000 --currency OLT
olclient send -c Alice --party Alice --counterparty Carol --amount 5099 --currency OLT
olclient send -c Bob --party Bob --counterparty Carol --amount 5012 --currency OLT
olclient send -c David --party David --counterparty Bob --amount 5000 --currency OLT
olclient send -c David --party David --counterparty Alice --amount 45000 --currency OLT
olclient send -c Bob --party Bob --counterparty Carol --amount 543 --currency OLT
olclient send -c David --party David --counterparty Carol --amount 5001 --currency OLT

sleep 8
