#!/bin/bash
#
# Install the OneLedger protocol and all of it's dependencies
#
# All output goes to the terminal and to a log
#
SETUP=$OLROOT/protocol/node/setup
echo "Installing OneLedger Dependencies" | tee $SETUP/install.log

#Install Bitcoin
echo "Installing Bitcoin" | tee -a $SETUP/install.log
$SETUP/install_bitcoin.sh | tee -a $SETUP/install.log
echo "====== Bitcoin finished =====" | tee -a $SETUP/install.log


#Install Ethereum
echo "Installing Ethereum" | tee -a $SETUP/install.log
$SETUP/install_ethereum.sh | tee -a $SETUP/install.log
echo "===== Ethereum finished =====" | tee -a $SETUP/install.log

##Install Consensus
#echo "Installing consensus" | tee -a $SETUP/install.log
#$SETUP/install_consensus.sh | tee -a $SETUP/install.log
#echo "===== Install finished =====" | tee -a $SETUP/install.log
