#!/bin/bash

if [ ! -d "$OLDATA/docker" ]
then
    mkdir "$OLDATA"/docker
fi

sudo rm -rf $OLDATA/docker/*-Node
sudo rm -rf $OLDATA/docker/bin/*

#start build container
docker-compose -f docker-compose-build.yml up -d

#wait until build container creates config files
until [ -n "$(find ~/oldata/docker -name '*Node' -type d)" ]; do
  echo "building..."
  sleep 5
done

#start node containers
docker-compose up -d