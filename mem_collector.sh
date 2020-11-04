#!/bin/bash

ps -o pid,user,%mem,command ax | sort -b -k3 -r | grep olfullnode >>/opt/data/devnet/ST_NOEVI_NOREWARDS/ram_usage.log
pmap 14797 | tail -n 1 >>/opt/data/devnet/ST_NOEVI_NOREWARDS/ram_usage.log
curl -sK -v http://localhost:45131/debug/pprof/heap >/opt/data/devnet/ST_NOEVI_NOREWARDS/$(date +\%Y\%m\%d\%H\%M\%S)-heap.out
echo $(date +\%Y\%m\%d\%H\%M\%S) >>/opt/data/devnet/ST_NOEVI_NOREWARDS/ram_usage.log
echo "-------------------------------------------------------------------------------------------------" >>/opt/data/devnet/ST_NOEVI_NOREWARDS/ram_usage.log
