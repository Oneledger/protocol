#!/usr/bin/env bash
VOL_DIR=$OLROOT/protocol/VOLUME
STAGING=$VOL_DIR/staging

if ! [ -f $VOL_DIR/David-Node/config/genesis.json ]
then
	tendermint testnet --v 4 \
		 --o $VOL_DIR/staging \
		 --populate-persistent-peers \
		 --starting-ip-address 192.167.11.1 \
		 --p2p-port 26611

	for node in David-Node#node0 Alice-Node#node1 Bob-Node#node2 Carol-Node#node3
	do
		nodeName=`echo $node | awk -F# '{print $1}'`
		dirName=`echo $node | awk -F# '{print $2}'`
		NODE_DIR=$VOL_DIR/$nodeName
		mkdir -p $NODE_DIR/tendermint/config
		touch $NODE_DIR/tendermint/config/genesis.json
		cat $STAGING/$dirName/config/genesis.json | jq -f $OLSCRIPT/genesis.jq > $NODE_DIR/tendermint/config/genesis.json
		rm $STAGING/$dirName/config/genesis.json
		cp -r $STAGING/$dirName/config $NODE_DIR/tendermint
	done
fi

rm -r $STAGING
