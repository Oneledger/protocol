#!/bin/bash

ps -o pid,user,%mem,command ax | sort -b -k3 -r | grep olfullnode >>/opt/ST_FULL_APP/ram_usage.log
curl -sK -v http://localhost:42637/debug/pprof/heap >/opt/ST_FULL_APP/$(date +\%Y\%m\%d\%H\%M\%S)-heap.out
echo $(date +\%Y\%m\%d\%H\%M\%S) >>/opt/ST_FULL_APP/ram_usage.log
echo "-------------------------------------------------------------------------------------------------" >>/opt/ST_FULL_APP/ram_usage.log
