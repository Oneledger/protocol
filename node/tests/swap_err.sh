#!/bin/bash

# The swaps hang up the node if the data is invalid
TESTS=$GOPATH/src/github.com/Oneledger/protocol/node/tests
$TESTS/resetStart.sh

olclient swap --root $OLDATA/Alice-Node --party Bob --counterparty Alice --amount 2.5
olclient swap --root $OLDATA/Alice-Node --party Alice --counterparty Bob --amount 2.5

olclient swap --root $OLDATA/Alice-Node --party Bob --counterparty Alice --amount 2.5
olclient swap --root $OLDATA/Alice-Node --party Alice --counterparty Bob --amount 2.5
