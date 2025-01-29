# Yesterday's News

The videos you are seeing are a random reshuffling of news clips from the day before today. The text is generated from the captions of those news clips, alongside a journal of mine from late October to early November 2016. The videos are purely random, but the text follows its own internal logic. 

Every day a new cycle, and we are meant to believe that _yesterday's news_ is no longer relevant. But every day leaves its mark, and finds its way back into the story as we continue forward. Each day follows another, and is inextricably linked to the one that came before it.

—[Andy](https://andydayton.com)

## Icons

 - [Fullscreen icon by Q.P. at the Noun Project](https://thenounproject.com/icon/fullscreen-6938590/)
 - [About icon by Mas Dhimas at the Noun Project](https://thenounproject.com/icon/about-6264304/)

## Contents

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

Server serves the actual video player, and generates captions using the model JSON file. The metadata and model are fetched from a remote Object Store by the server, and video clips from the Object Store are loaded via a CDN. 

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

Builder will download videos from YouTube, cut them up, build a markov model, and upload everything to an S3-compatible Object Store.

Dependencies:

 - Go 1.22.10
 - `yt-dlp` (version `2024.12.06` or higher)
 - `ffmpeg`
 - `ffprobe`

### Usage

To get usage info, run `go run ./app/builder/main.go -h`.

Run the default build and save artifacts:

```bash
make builder
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
	BuildPhaseAll            BuildPhase = "all"
	BuildPhaseSetup          BuildPhase = "setup"
	BuildPhaseDownloadVideos BuildPhase = "download-videos"
	BuildPhaseCutVideos      BuildPhase = "cut-videos"
	BuildPhaseUploadVideos   BuildPhase = "upload-videos"
	BuildPhaseModel          BuildPhase = "model"
	BuildPhaseManifest       BuildPhase = "manifest"
	BuildPhasePromote        BuildPhase = "promote"
	BuildPhaseCleanup        BuildPhase = "cleanup"
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
./app/builder/deploy.sh <server IP>
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
