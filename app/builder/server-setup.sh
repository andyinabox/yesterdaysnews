#!/bin/bash

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <serverhost>"
    exit 1
fi

HOST=$1

scp .env yesterdaysnews@$HOST:.env
scp bin/builder-linux-amd64 yesterdaysnews@$HOST:builder
