#!/bin/bash

sudo rm -rf $OLDATA/docker/*-Node

#start build container
docker-compose -f docker-compose-build.yml up -d

#wait until build container creates config files
until [ -n "$(find ~/oldata/docker -name '*Node' -type d)" ]; do
  echo "building..."
  sleep 3
done

#start node containers
docker-compose up -d