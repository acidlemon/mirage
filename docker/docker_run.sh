#!/bin/sh

docker run -rm -v /mirage -v /var/run/docker.sock:/var/run/docker.sock  -t mirage:latest
