#!/bin/bash

docker build -t mapreduce-master -f Dockerfile.master .
docker build -t mapreduce-worker -f Dockerfile.worker .
