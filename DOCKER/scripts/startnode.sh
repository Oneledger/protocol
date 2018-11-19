#!/usr/bin/env bash

# Get the node id for each fullnode
function generate_peers()
{
	local addrLength=40
	local peers=()
	for peer in $(tr "," " " <<< $1); do
		local rpc_peer="$(echo $peer | awk -F: '{print $1}'):$OL_PORT_RPC"
		local peerID="$(curl -s $rpc_peer/status | jq ".result.node_info.id" | tr -d '\"')"
		if [ "${#peerID}" == "$addrLength" ]; then
			peers+=("$peerID@$peer")
		fi
	done
	echo `tr " " "," <<< "${peers[*]}"`
}

function dl_genesis()
{
	local genesis
	local peer_address=$1
	local rpc_peer="$(echo $peer_address | awk -F: '{print $1}'):$OL_PORT_RPC"
	local genesis=$(curl -s ${rpc_peer}/genesis | jq .result.genesis)
	if [ -z "$genesis" ]; then
		echo "Genesis download failed"
		echo "RPC Peer ${rpc_peer} looks unhealthy"
		exit 1
	else
		echo $genesis
	fi
	local status=$?
	exit $status
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
olLog=$LOG/olfullnode.log
tmData=$LOG/tendermint


# If we are a fullnode joining an existing network, download the genesis
# file and get the right node id from each peer ip listed from OL_PEERS
if [ "$OL_PEERS" ]; then
	echo
	echo "Making directory $LOG"
	echo
	mkdir -p $tmData/config
	touch $tmLog
	touch $olLog


	foundPeers=$(generate_peers $OL_PEERS)
	peer1=$(echo $OL_PEERS | awk -F, '{print $1}')
	genesis_file=$(dl_genesis $peer1)
	echo "Genesis: $genesis_file"
	if [ "$?" != "0" ]; then
		echo "Genesis download failed"
		exit 1
	fi

	echo "Copying genesis file..."
	echo $genesis_file >> $tmData/config/genesis.json
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

olfullnode node \
	--root $OLDATA/$nodeName/olfullnode \
	--node $nodeName \
	--app $appAddress \
	--address $rpcAddress \
	>> $olLog 2>&1 &
	# TODO: Add btc and eth rpc addresses

tail -f $tmLog
