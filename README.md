# Yesterday's News

This is really two projects in one

 - A Go port of [yesterdays-news-py](https://github.com/andyinabox/yesterdays-news-py/). The goal is to build assets and them push them to an object store.
 - A web project that will be like [yesterdays-news-of](https://github.com/andyinabox/yesterdays-news-of/) in a browser.


## Directory structure

 - `app` - entrypoints and data for the two main applications 
 - `bin` - compiled binaries end up here
 - `cmd` - misc utility command entrypoints
 - `dist` - this is where all downloaded and generated files end up 
 - `domain` - code specific to this project
 - `pkg` - more general-use code that could be used in other projects
 - `test` - test fixtures

## Server

Dependencies:

 - Go 1.22.10
 - `docker`

### Usage

To get usage info, run `go run ./app/server/main.go -h`.

To run the server locally with good defaults for local development:

```bash
make server
```

This will:
 - Enable loading assets from the filesystem so you can do frontend development
 - Set to check for updated build assets (on the Object Store) to every 20 seconds
 - Enable verbose logging


### Building the server

The server is deployed using docker:

```bash
# build the server conainer
make docker-build-server
# test the server container
make docker-run-server
# push the server container to docker hub
make docker-push-server
```

## Builder

Builder will download videos from YouTube, cut them up, build a markov model, and upload everything to an S3-compatible Object Store. Some helpful commands:

Dependencies:

 - Go 1.22.10
 - `yt-dlp` (version `2024.12.06` or higher)
 - `ffmpeg`
 - `ffprobe`

### Usage

To get usage info, run `go run ./app/builder/main.go -h`.

Run the build and save artifacts:

```bash
make builder
# or 
go run ./app/builder/main.go --keepoutput
```

Test the build without uploading artifacts (helpful if you want to populate "dist" for use with `objectstoremock`):

```bash
go run ./app/builder/main.go --keepoutput --skipupload
```

Test individual build steps:

```bash
go run ./app/builder/main.go -b <build step>
```

You can find the proper name for each build step in [domain/builder.go](domain/builder.go).

```go
const (
	BuildPhaseAll        BuildPhase = "all"
	BuildPhaseSetup      BuildPhase = "setup"
	BuildPhaseVideoClips BuildPhase = "video-clips"
	BuildPhaseModel      BuildPhase = "model"
	BuildPhaseManifest   BuildPhase = "manifest"
	BuildPhasePromote    BuildPhase = "promote"
	BuildPhaseCleanup    BuildPhase = "cleanup"
)
```

### Builder remote server setup

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


## Utils

### `objectstoremock`

Serves a local Object Store using the contents of `dist`. 

To get usage info, run `go run ./cmd/objectstoremock/main.go -h`.

Example usage:

```bash
make objectstoremock
# or
./cmd/objectstoremock/main.go
```


### `getplaylistid`

Gets the ID for the main playlist of a channel using the channel handle (i.e. `@CNN`).

To get usage info, run `go run ./cmd/getplaylistid/main.go -h`.

Example usage:

```bash
./cmd/getplaylistid/main.go -n <channel handle>
```
