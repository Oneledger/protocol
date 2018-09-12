#!/usr/bin/env bash

nodeName=${ID}-Node

OL_BTC_ADDRESS=NONE
OL_ETH_ADDRESS=NONE

prefix="tcp://0.0.0.0:"

rpcAddress="$prefix$OL_PORT_RPC"
p2pAddress="$prefix$OL_PORT_P2P"
appAddress="$prefix$OL_PORT_APP"

LOG=$OLDATA/$nodeName
tmLog=$LOG/tendermint.log
olLog=$LOG/fullnode.log
tmData=$LOG/tendermint

echo "============================================================" >> $tmLog
echo "Starting Tendermint" >> $tmLog
echo "============================================================" >> $tmLog

tendermint node --home $LOG/tendermint \
	--moniker $nodeName \
	--rpc.laddr $rpcAddress\
	--p2p.laddr $p2pAddress \
	--home $tmData \
	--proxy_app $appAddress \
	>> $tmLog 2>&1 &

echo "============================================================" >> $olLog
echo "Starting Fullnode" >> $olLog
echo "============================================================" >> $olLog

fullnode node \
	--root $OLDATA/$nodeName/fullnode \
	--node $nodeName \
	--app $appAddress \
	--address $rpcAddress \
	--debug \
	>> $olLog 2>&1 &
	# TODO: Add btc and eth rpc addresses

tail -f $olLog
#tail -f $tmLog
