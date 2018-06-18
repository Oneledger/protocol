#!/bin/bash
#
# Install the OneLedger protocol and all of it's dependencies
#

#Install Bitcoin
echo "Installing Bitcoin"
sudo ./install_bitcoin.sh
echo "======Bitcoin finished====="


#Install Ethereum
echo "Installing Ethereum"
./install_ethereum.sh
echo "=====Ethereum finished====="

#Install Consensus
echo "Installing consensus"
./install_consensus.sh
echo "=====Install finished====="



