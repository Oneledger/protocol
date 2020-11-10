#!/bin/bash

#Stop nodes
# shellcheck disable=SC2046
docker container kill $(docker ps -q)

#Build Binaries
docker start builder
docker exec -it builder bash -c "cd /home/ubuntu/go/protocol && make install_c"
docker stop builder

#Start nodes
docker-compose up -d