#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

$CMD/startOneLedger

# assumes olfullnode is in the PATH
SEQ=`$CMD/nextSeq`
olclient send --root $OLDATA/Bob-Node --party Bob --counterparty Alice --amount 1000.1 --currency OLT --fee 0.01
olclient send --root $OLDATA/Alice-Node --party Alice --counterparty David --amount 5001.2 --currency OLT --fee 0.02
olclient send --root $OLDATA/Bob-Node --party Bob --counterparty David --amount 523.3 --currency OLT --fee 0.03
olclient send --root $OLDATA/Bob-Node --party Bob --counterparty David --amount 1000.4 --currency OLT --fee 0.01
olclient send --root $OLDATA/Alice-Node --party Alice --counterparty Carol --amount 2000.5 --currency OLT --fee 0.04
olclient send --root $OLDATA/Alice-Node --party Alice --counterparty David --amount 5099.6 --currency OLT --fee 0.05
olclient send --root $OLDATA/David-Node --party David --counterparty Carol --amount 51.7 --currency OLT --fee 0.06
olclient send --root $OLDATA/David-Node --party David --counterparty Bob --amount 230.8 --currency OLT --fee 0.0101
olclient send --root $OLDATA/Alice-Node --party Alice --counterparty Bob --amount 5050.9 --currency OLT --fee 0.0102
olclient send --root $OLDATA/Bob-Node --party Bob --counterparty Alice --amount 50200.01 --currency OLT --fee 0.03
olclient send --root $OLDATA/Alice-Node --party Alice --counterparty Bob --amount 3050.02 --currency OLT --fee 0.0203
olclient send --root $OLDATA/David-Node --party David --counterparty Alice --amount 10002.03 --currency OLT --fee 0.01
olclient send --root $OLDATA/Alice-Node --party Alice --counterparty Carol --amount 5099.04 --currency OLT --fee 0.04
olclient send --root $OLDATA/Bob-Node --party Bob --counterparty Carol --amount 2012.05 --currency OLT --fee 0.07
olclient send --root $OLDATA/David-Node --party David --counterparty Bob --amount 5000.06 --currency OLT --fee 0.120
olclient send --root $OLDATA/David-Node --party David --counterparty Alice --amount 45000.07 --currency OLT --fee 0.0301
olclient send --root $OLDATA/Bob-Node --party Bob --counterparty Carol --amount 543.08 --currency OLT --fee 0.07
olclient send --root $OLDATA/David-Node --party David --counterparty Carol --amount 5001.09 --currency OLT --fee 0.02

sleep 8
