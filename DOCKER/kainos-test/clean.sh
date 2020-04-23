#!/bin/bash

pkill olfullnode

rm *.log
rm -rf nodedata/*
rm -rf consensus/data/*.db
rm -rf consensus/data/cs.wal
rm -rf consensus/config/write*
rm -rf consensus/config/genesis.json
wget https://raw.githubusercontent.com/Oneledger/protocol/develop/DOCKER/kainos-test/genesis.json

mv genesis.json ./consensus/config/
echo '{
  "height": "0",
  "round": "0",
  "step": 0
}'  > ./consensus/data/priv_validator_state.json

