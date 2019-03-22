#!/bin/bash

if [[ "$OSTYPE" == "linux-gnu" ]]; then
	sudo apt install libleveldb-dev libsnappy-dev
	echo
	echo "........."
	echo "Done!"
else
	echo "Your OS isn't supported, install leveldb and snappy libs on your own"
fi
