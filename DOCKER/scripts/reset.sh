#!/usr/bin/env bash

VOL_DIR="$OLROOT/protocol/VOLUME"

if [ -d "$VOL_DIR" ]
then
	echo "Clearing VOLUME directory..."
	rm -r $VOL_DIR
fi
