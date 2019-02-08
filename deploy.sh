#!/bin/bash

curl https://ipinfo.io/ip >> IPV4

GOPATH="/home/ubuntu"
echo export GOPATH="/home/ubuntu" >> ~/.bashrc
echo export PATH="$PATH:$GOPATH/bin" >> ~/.bashrc
source ~/.bashrc
chmod 555 -R "$GOPATH/bin"  
#work around. Need to modify CLI for olfullnode to accept full directory path as input for co$
cp "$GOPATH/config.toml" /tmp 
GOPATH="~" PATH="$PATH:$GOPATH/bin" olfullnode init devnet
