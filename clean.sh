#!/usr/bin/env bash
docker stop $(docker ps -q)

docker rm $(docker ps -aq)
