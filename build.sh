#!/bin/bash

docker build -t mapreduce-master -f Dockerfiles/Dockerfile.master .
docker build -t mapreduce-worker -f Dockerfiles/Dockerfile.worker .
