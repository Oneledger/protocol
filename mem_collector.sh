#!/bin/bash

ps -o pid,user,%mem,vsz,rss,command ax | sort -b -k3 -r | grep olfullnode >>/opt/ST_COMM/ram_usage.log
curl -sK -v http://localhost:33077/debug/pprof/heap >/opt/ST_COMM/$(date +\%Y\%m\%d\%H\%M\%S)-heap.out
echo $(date +\%Y\%m\%d\%H\%M\%S) >>/opt/ST_COMM/ram_usage.log
echo "-------------------------------------------------------------------------------------------------" >>/opt/ST_COMM/ram_usage.log
