#!/bin/bash

if [ ! -d "${OLDATA/olfullnode}" ]; then
	mkdir -p "${OLDATA/olfullnode}"
	touch ${OLDATA}/olfullnode.log
fi

if [ ! -d "${OLDATA/consensus}" ]; then
	mkdir -p "${OLDATA/consensus}"
	touch ${OLDATA}/consensus.log
fi

olfullnode node -c config --root ${OLDATA}/olfullnode --tendermintRoot ${OLDATA}/consensus \
      --seeds 6f724aaeb0a9232b91060fc98af9f034fa7970d0@35.203.120.202:26615,f2815fcb2deb6209c26d86e1cbb6f7708f21e458@35.203.59.66:26615 --persistent_peers f2815fcb2deb6209c26d86e1cbb6f7708f21e458@35.203.59.66:26615,6f724aaeb0a9232b91060fc98af9f034fa7970d0@35.203.120.202:26615 \
       >> ${OLDATA}/olfullnode.log 2>&1 &

sleep 2

