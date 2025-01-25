# Yesterday's News

This is really two projects in one

 - A Go port of [yesterdays-news-py](https://github.com/andyinabox/yesterdays-news-py/). The goal is to build assets and them push them to an object store.
 - A web project that will be like [yesterdays-news-of](https://github.com/andyinabox/yesterdays-news-of/) in a browser.


Dependencies:

 - Go 1.22.10
 - `yt-dlp` 2024.12.06
 - `ffmpeg`
 - `ffprobe`
 - `docker` for publishing the webserver

## Directory structure

 - `app` - entrypoints and data for the two main applications 
 - `bin` - compiled binaries end up here
 - `cmd` - misc utility command entrypoints
 - `dist` - this is where all downloaded and generated files end up 
 - `domain` - code specific to this project
 - `pkg` - more general-use code that could be used in other projects
 - `test` - test fixtures

## TODO

Todo list can be found in the [GitLab Issues](https://gitlab.com/andyinabox/yesterdaysnews/-/issues) currently


## Server

To get usage info, run `go run ./app/server/main.go -h`. For example:

```
Usage of server:
  -a	load assets from filesystem (for easier frontend development)
  -m string
    	manifest check interval (default "1h")
  -maxd float
    	max caption delay in seconds (default 5)
  -maxl int
    	max caption length in words (default 15)
  -mind float
    	min caption delay in seconds (default 1.5)
  -minl int
    	min caption length in words (default 5)
  -p int
    	markov chain prefix length (default 2)
  -port int
    	server port (default 8080)
  -v	verbose logging
```

## Builder

Builder will download videos from YouTube, cut them up, build a markov model, and upload everything to an S3-compatible Object Store.

### Server setup

I am currently running this on a server with the following attributes:

 - Linux Ubuntu 24.04 LTS 64-bit
 - 16GB RAM
 - 50GB Disk
 - 4 CPUs
 - Cloud-init: [cloud-config.yml](app/builder/cloud-config.yml)

I initially tried running it on a smallar instance but found it froze up.

After you have provisioned the server, grab the IPv4 and run

```
./app/builder/setup.sh <server IP>
```

This will copy the remaining necessary files to the server.

### Usage

To get usage info, run `go run ./app/builder/main.go -h`. For example:

```
Usage of builder:
  -b string
    	build phase to execute (default "all")
  -buildstokeep int
    	total completed builds to keep when cleaning up (default 5)
  -containername string
    	object storage container name (default "yesterdaysnews")
  -count int
    	download count per playlist (default 10)
  -hospitalweight int
    	weight for the hospital corpus in chain (default 1)
  -keepoutput
    	keep artifacts after successful build
  -maxcliplength int
    	maximum clip length in seconds (default 15)
  -maxvideosize int
    	max video download size in bytes (default 52428800)
  -mincliplength int
    	minimum clip length in seconds (default 5)
  -newsweight int
    	weight for the news corpus in chain (default 1)
  -output string
    	dir to output build artifacts to (default "dist")
  -playlistids string
    	comma-separated list of playlists to download (default "UUupvZG-5ko_eiXAupbDfxWw,UUaXkIU1QidjPwiAYu6GcHjg,UUXIJgqnII2ZOINSWNOGFThA")
  -prefixlength int
    	caption chain prefix length (default 2)
  -skipupload
    	skip upload step for builds
  -v	verbose output
```

## Utils

### `objectstoremock`

Serves a local Object Store using the contents of `dist`. 

To get usage info, run `go run ./cmd/objectstoremock/main.go -h`. For example:

```
Usage of objectstoremock:
  -d string
    	dir to serve (default "dist")
  -p int
    	port to serve on (default 9000)
  -v	verbose logging
```

### `getplaylistid`

Gets the ID for the main playlist of a channel using the channel handle (i.e. `@CNN`).

To get usage info, run `go run ./cmd/getplaylistid/main.go -h`. For example:

```
Usage of getplaylistid:
  -n string
    	channel name
  -v	verbose output
```