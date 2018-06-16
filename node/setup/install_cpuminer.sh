#! /bin/bash

apt-get update -qq && \
apt-get install -qqy automake libcurl4-openssl-dev git make

git clone https://github.com/pooler/cpuminer

cd cpuminer && \
./autogen.sh && \
./configure CFLAGS="-O3" && \
make

ln -s ./minerd /usr/bin/minerd