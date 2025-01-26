#!/bin/bash

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <serverhost>"
    exit 1
fi

HOST=$1
USER=yesterdaysnews
REMOTE_ENV_PATH=/home/$USER/.env
REMOTE_BIN_PATH=/home/$USER/builder
REMOTE_SCRIPT_PATH=/home/$USER/run.sh

# ensure that a fresh binary is there
echo "building fresh builder binary..."
rm bin/builder-linux-amd64
make bin/builder-linux-amd64

# check for .env file
if ! ssh $USER@$HOST [ -f "$REMOTE_ENV_PATH" ]; then
    echo ".env file not found, copying from local..."
    scp .env $USER@$HOST:$REMOTE_ENV_PATH
fi

# check for run.sh file
if ! ssh $USER@$HOST [ -f "$REMOTE_SCRIPT_PATH" ]; then
    echo "run.sh not found, copying from local..."
    scp app/builder/run.sh $USER@$HOST:$REMOTE_SCRIPT_PATH
fi

echo "deploying latest builder binary..."
scp bin/builder-linux-amd64 $USER@$HOST:$REMOTE_BIN_PATH

echo "done."