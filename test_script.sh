#!/bin/bash

pwd

# to build binary
make install_c

# unit testing
make utest

# apply validator test
make applytest

# ONS test, oneledger naming service tests
make onstest

# withdraw test
make withdrawtest

# governance test
make govtest

# all tests
make alltest

# make rpc auth test
make rpcAuthtest

# generate coverage result, cover.html under protocol
make coverage
