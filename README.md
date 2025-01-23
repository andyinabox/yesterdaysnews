# Yesterday's News

This is really two projects in one

 - A Go port of [yesterdays-news-py](https://github.com/andyinabox/yesterdays-news-py/). The goal is to build assets and them push them to an object store.
 - A web project that will be like [yesterdays-news-of](https://github.com/andyinabox/yesterdays-news-of/) in a browser.

## Todo

Todo list can be found in the [GitLab Issues](https://gitlab.com/andyinabox/yesterdaysnews/-/issues) currently

## Builder

## Server

Below is a command dump to track what I've done to the builder server:

```bash

# ------------------------------------------------------
# via https://docs.docker.com/engine/install/ubuntu/
# ------------------------------------------------------

sudo apt-get update
sudo apt-get install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources:
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update

sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# ------------------------------------------------------
# freestyling the rest below
# ------------------------------------------------------

mkdir dist
vim .env
# added env variables


```

