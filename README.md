# Yesterday's News

The videos you are seeing are news clips from the day before today.

The text is generated from the subtitles for those videos, alongside a journal of mine from late October to early November 2016.

The video ordering is purely random, but the text follows its own internal logic.

—[Andy](https://andydayton.com)

## Icons

 - [Fullscreen icon by Q.P. at the Noun Project](https://thenounproject.com/icon/fullscreen-6938590/)
 - [About icon by Mas Dhimas at the Noun Project](https://thenounproject.com/icon/about-6264304/)

## Infrastructure

I've made an effort to use cloud infratructure that is based in Europe and run somewhat environmentally sustainably. I've currently settled on three different services in order to keep costs relatively low:

 - The Server application is hosted on [Infomaniak](https://www.infomaniak.com/)
 - The Builder application is run on [Scaleway](https://www.scaleway.com/en/)
 - Object Store assets stored on [Exoscale](https://www.exoscale.com/)

## Directory structure

 - `app` - entrypoints and data for the two main applications 
 - `bin` - compiled binaries end up here
 - `cmd` - misc utility command entrypoints
 - `dist` - this is where all downloaded and generated files end up 
 - `domain` - code specific to this project
 - `pkg` - more general-use code that could be used in other projects
 - `test` - test fixtures

## Primary Applications

### Server (`app/server`)

Server serves the actual video player, and generates captions using the model JSON file. The metadata and model are fetched from a remote Object Store by the server, and video clips from the Object Store are loaded via a CDN. 

Dependencies:

 - Go 1.22.10
 - `docker`


### Builder (`app/builder`)

Builder will download videos from YouTube, cut them up, build a markov model, and upload everything to an S3-compatible Object Store.

Dependencies:

 - Go 1.22.10
 - `yt-dlp` (version `2024.12.06` or higher)
 - `ffmpeg`
 - `ffprobe`

Additionaally the builder requires more resources to work well, so it's a good idea to give it more RAM and CPUs.

## Running locally

First you will want to run the Builder locally:

```bash
make builder-local
```

This will download all the assets and process them, but skip the step of uploading them to the Object Store.

Next you'll want to start a local Object Store mock to serve the assets:

```bash
make objectstoremock
```

This will similate the S3-compatible Object Store used provide assets in production, but using the locally processed assets.

Finally, **in a new terminal window** (you need to keep the `objectstoremock` running), run the following:

```bash
make server-local
```

This will run the server on `localhost`, retrieving assets from the `objectstoremock`.

## Releasing

There is a script `release.sh` that automates tagging, building the docker container, and publishing to the container registry. It is run separately for each application:

```bash
# create a v0.0.0 release for the server application
./release.sh server v0.0.0

# create a v0.0.0 release for the builder application
./release.sh builder v0.0.0
```


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
