#!/bin/bash

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <serverhost>"
    exit 1
fi

HOST=$1

# ensure that binary is there
make bin/builder-linux-amd64

# copy necessay files over
scp .env yesterdaysnews@$HOST:.env
scp bin/builder-linux-amd64 yesterdaysnews@$HOST:builder
