#!/bin/bash


export BTCDATA=$OLDATA/bitcoin
export LOG=$OLDATA

echo "Killing all bitcoin process" >> $LOG/bitcoin.log

killall --regex bitcoin.*

echo "Running Bitcoin process: $(pgrep bitcoin.* | wc -l)" >> $LOG/bitcoin.log