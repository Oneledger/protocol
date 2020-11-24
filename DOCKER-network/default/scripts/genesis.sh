#!/bin/bash

echo "Deleting Node Data"
for name in $(ls -l $OLDATA | grep Node  | awk '(NR>0){print $9}')
do
    rm -rf $OLDATA/$name/nodedata/*
    rm -rf $OLDATA/$name/consensus/data/*.db
    rm -rf $OLDATA/$name/consensus/data/*.wal
    rm -rf $OLDATA/$name/consensus/config/addrbook.json

    echo "{
  \"height\": \"0\",
  \"round\": \"0\",
  \"step\": 0
}" > $OLDATA/$name/consensus/data/priv_validator_state.json

    #Copy new genesis file to node folders
    cp -f $OLDATA/genesis.json $OLDATA/$name/consensus/config/
done
