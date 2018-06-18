#!/bin/sh
#
# Install the OneLedger protocol and all of it's dependencies
#

#Install Bitcoin
./install_bitcoin.sh

#Install Ethereum
./install_etheruem.sh
./install_cpuminer.sh

#Install Consensus
./install_consensus.sh

#Install Protocol
./install_prototype.sh

