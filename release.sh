#!/bin/bash

set -e

USAGE="Usage: $0 <builder/server> <tag>"

if [ "$#" -ne 2 ]; then
  echo $USAGE
  exit 1
fi

if test "$1" != "builder" && test "$1" != "server"; then
  echo $USAGE
  exit 1
fi

APP=$1
TAG=$2

# build a fresh binary
make clean-bin bin/$APP-linux-amd64

# build the container
docker buildx build --platform linux/arm64\
 -f app/$APP/Dockerfile\
 -t andyinabox/yesterdaysnews-$APP:$TAG .

# push tagged container
docker push andyinabox/yesterdaysnews-$APP:$TAG

# tag as latest and push
docker tag andyinabox/yesterdaysnews-$APP:$TAG andyinabox/yesterdaysnews-$APP:latest
docker push andyinabox/yesterdaysnews-$APP:latest
