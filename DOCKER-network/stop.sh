#!/bin/bash

docker container kill $(docker ps -q)
docker system prune -f