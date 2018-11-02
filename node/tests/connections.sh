#!/bin/bash
#
# This sequence of commands seems to cause the olfullnodes to fail if the CPU is low
#

echo "========== Resetting =========="
#make reset

echo "========== Test =========="
pushd ..
make test
sleep 6

echo "========== Swaptest =========="
make swaptest
popd
sleep 6

echo "========== Register =========="
./register.sh
sleep 6

echo "========== Swap =========="
./swap.sh
sleep 6

echo "========== Identity =========="
./identity.sh
sleep 6

echo "========== Account =========="
./accounts.sh
sleep 6

