#!/usr/bin/env bash

# Get the node id for each fullnode
function generate_peers()
{
	addrLength=40
	peers=()
	for peer in $(tr "," " " <<< $1); do
		rpcPeer="$(echo $peer | awk -F: '{print $1}'):$OL_PORT_RPC"
		peerID="$(curl -s $rpcPeer/status | jq ".result.node_info.id" | tr -d '\"')"
		echo $peerID > $LOG/2peers
		if [ "${#peerID}" == "$addrLength" ]; then
			peers+=("$peerID@$peer")
		fi
	done
	echo `tr " " "," <<< "${peers[*]}"`
}

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

## If the PEERS argument exists, get the genesis block from specified peers
#if [ $OL_PEERS ] && ! [ -f "$tmData/config/genesis.json" ]
#then
#	mkdir -p "$tmData/config"
#	echo
#	echo "Downloading genesis block from peers..."
#	try_peer_genesis $OL_PEERS
#
##	echo
##	echo "Copying genesis block to node..."
##	$genesis > $tmData/config/genesis.json
#fi

if [ "$OL_PEERS" ]; then
	touch $tmLog
	touch $olLog
	foundPeers=$(generate_peers $OL_PEERS)
	echo
	echo "Starting node with peers: $foundPeers"
	echo
	peerArgs="--p2p.persistent_peers=\"$foundPeers\""
else
	peerArgs=""
fi

echo "============================================================" >> $tmLog
echo "Starting Tendermint" >> $tmLog
echo "============================================================" >> $tmLog

tendermint node --home $LOG/tendermint \
	--moniker $nodeName \
	--rpc.laddr $rpcAddress\
	--p2p.laddr $p2pAddress \
	--home $tmData \
	--proxy_app $appAddress \
	${peerArgs:-""} >> $tmLog 2>&1 &

echo "============================================================" >> $olLog
echo "Starting Fullnode" >> $olLog
echo "============================================================" >> $olLog

fullnode node \
	--root $OLDATA/$nodeName/fullnode \
	--node $nodeName \
	--app $appAddress \
	--address $rpcAddress \
	>> $olLog 2>&1 &
	# TODO: Add btc and eth rpc addresses

tail -f $tmLog
