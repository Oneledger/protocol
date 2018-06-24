#!/bin/bash

go get -u github.com/tendermint/tendermint/cmd/tendermint
#go get -u github.com/tendermint/abci/cmd/abci-cli

pushd $GOPATH/src/github.com/tendermint/tendermint
git checkout tags/v0.18.0
make get_tools
make get_vendor_deps
make install
#make test

