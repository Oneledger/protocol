#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

$CMD/startOneLedger

# assumes olfullnode is in the PATH
SEQ=`$CMD/nextSeq`
olclient send --root $OLDATA/Bob-Node --party Bob --counterparty Alice --amount 10.1 --currency OLT --fee 0.101
olclient send --root $OLDATA/Alice-Node --party Alice --counterparty David --amount 501.2 --currency OLT --fee 0.102
olclient send --root $OLDATA/Bob-Node --party Bob --counterparty David --amount 53.3 --currency OLT --fee 0.103
olclient send --root $OLDATA/Bob-Node --party Bob --counterparty David --amount 10.4 --currency OLT --fee 0.101
olclient send --root $OLDATA/Alice-Node --party Alice --counterparty Carol --amount 20.5 --currency OLT --fee 0.104
olclient send --root $OLDATA/Alice-Node --party Alice --counterparty David --amount 599.6 --currency OLT --fee 0.105
olclient send --root $OLDATA/David-Node --party David --counterparty Carol --amount 51.7 --currency OLT --fee 0.106
olclient send --root $OLDATA/David-Node --party David --counterparty Bob --amount 2.8 --currency OLT --fee 0.10101
olclient send --root $OLDATA/Alice-Node --party Alice --counterparty Bob --amount 55.9 --currency OLT --fee 0.10102
olclient send --root $OLDATA/Bob-Node --party Bob --counterparty Alice --amount 52.01 --currency OLT --fee 0.103
olclient send --root $OLDATA/Alice-Node --party Alice --counterparty Bob --amount 305.02 --currency OLT --fee 0.10203
olclient send --root $OLDATA/David-Node --party David --counterparty Alice --amount 102.03 --currency OLT --fee 0.201
olclient send --root $OLDATA/Alice-Node --party Alice --counterparty Carol --amount 599.04 --currency OLT --fee 0.204
olclient send --root $OLDATA/Bob-Node --party Bob --counterparty Carol --amount 212.05 --currency OLT --fee 0.207
olclient send --root $OLDATA/David-Node --party David --counterparty Bob --amount 50.06 --currency OLT --fee 0.120
olclient send --root $OLDATA/David-Node --party David --counterparty Alice --amount 450.07 --currency OLT --fee 0.20301
olclient send --root $OLDATA/Bob-Node --party Bob --counterparty Carol --amount 5.08 --currency OLT --fee 0.207
olclient send --root $OLDATA/David-Node --party David --counterparty Carol --amount 5.09 --currency OLT --fee 0.202

#sleep 8
