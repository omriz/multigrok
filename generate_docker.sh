#!/usr/bin/env bash

# An example script on how to build a docker file for the multigro binary.

set -x
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
docker build -t $1 .