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
GITTAG=$APP-$TAG

echo ""
echo "repo status:"
echo ""
git status --short --branch

echo ""
echo "existing tags:"
echo ""
git tag | grep $APP

echo ""
echo "App:     $APP"
echo "Tag:     $TAG"
echo "Git tag: $GITTAG"

echo ""
echo "Please review output above. Continue with tag/release workflow? (y/N)"
read CONFIRM

if [[ "$CONFIRM" != "y" ]]; then
    exit 0;
fi

git tag $GITTAG
git push --tags
git checkout $GITTAG

# build a fresh binary
make clean-bin bin/$APP-linux-amd64

# build the container
docker buildx build --platform linux/amd64\
 -f app/$APP/Dockerfile\
 -t andyinabox/yesterdaysnews-$APP:$TAG .

# push tagged container
docker push andyinabox/yesterdaysnews-$APP:$TAG

# tag as latest and push
docker tag andyinabox/yesterdaysnews-$APP:$TAG andyinabox/yesterdaysnews-$APP:latest
docker push andyinabox/yesterdaysnews-$APP:latest

git checkout main