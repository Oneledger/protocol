#!/bin/bash

# shellcheck disable=SC2046
docker container kill $(docker ps -q)
docker system prune -f