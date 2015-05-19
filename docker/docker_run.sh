#!/bin/sh

docker run --rm -v /data -v /var/run/docker.sock:/docker.sock -t mirage:latest
