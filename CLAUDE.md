# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run Commands

```bash
# Run tests
go test ./...
go test ./domain/videoprocessor    # single package

# Local development
make server-local                  # server with local mock object store (start objectstoremock first)
make objectstoremock               # local S3 mock server
make builder-local                 # builder without upload

# Production-like
make server                        # server with real object store credentials from .env
make builder                       # builder with upload

# Docker
make server-docker                 # build & run server container
make builder-docker                # build & run builder container

# Clean
make clean                         # remove all build artifacts (dist/, bin/, .assets/)
```

## Architecture

Two applications sharing a common domain layer:

- **Server** (`app/server/`) — Web app serving video clips and Markov-chain-generated captions. Loads a manifest from object storage, serves HTML templates with embedded web components (Lit/Haunted).
- **Builder** (`app/builder/`) — Batch job that downloads YouTube videos (via yt-dlp), cuts them into clips (via ffmpeg), builds a Markov chain model from subtitles, uploads everything to S3-compatible storage, and generates a manifest.

### Layer structure

- `domain/` — Business logic. Contains interfaces at the root (`domain/builder.go`, `domain/server.go`, `domain/video-processor.go`, etc.) and implementations in sub-packages (`domain/builder/`, `domain/server/`, `domain/videoprocessor/`, `domain/youtubeservice/`, `domain/containerservice/`, etc.).
- `pkg/` — Reusable, project-independent packages: `shell/` (command execution), `markov/` (chain library), `objectstoreclient/` (S3 wrapper), `youtubedownloader/` (yt-dlp wrapper), `srt/` (subtitle parser), `mediatool/` (ffmpeg/ffprobe wrapper), etc.
- `app/` — Entry points and application-specific assets.
- `cmd/` — CLI utilities (objectstoremock, esbuild wrapper, playlist tools).

### Builder phases (run individually or all together)

setup → download-videos → cut-videos → upload-videos → extract-images → poster-image → model → combine → manifest → promote → cleanup

### Frontend asset pipeline

Source in `app/server/assets/` → bundled via esbuild (`cmd/esbuild/`) → output to `app/server/.assets/` → embedded into server binary via `//go:embed`.

## Key patterns

- **Go module path**: `gitlab.com/andyinabox/yesterdaysnews`
- **Go version**: 1.22 (toolchain 1.22.10)
- **Dependency injection**: Services accept config structs and implement domain interfaces
- **Streaming**: Operations use channels for async progress/error reporting
- **Error handling**: Errors must always be handled, never discarded with `_`. Use `fmt.Errorf("...%w", err)` for wrapping. For non-fatal errors in streaming/batch operations, send errors to the `ErrorHandler` via `eh.Add(typ, err)` or `eh.Channel()` rather than returning early — this allows the operation to continue while tracking error counts per type and panicking if a threshold is exceeded.
- **Logging**: `github.com/charmbracelet/log` (structured, supports JSON format)
- **Config**: Environment variables loaded from `.env` via godotenv + `codingconcepts/env` tags

## External tool dependencies (builder only)

- `yt-dlp` — YouTube video downloading
- `ffmpeg` / `ffprobe` — Video processing
- Paths configurable via `YN_YT_DLP_PATH`, `YN_FFMPEG_PATH`, `YN_FFPROBE_PATH` env vars

## Release

```bash
./release.sh server v0.1.0              # tag, build, push server image
./release.sh builder v0.1.0             # tag, build, push builder image
./release.sh server v0.1.0 --dry-run    # print every mutating step without executing it
```

Publishes to the Forgejo container registry at `code.andydayton.com` as
`code.andydayton.com/andy/yesterdaysnews-{server|builder}:{tag}` (also tagged `:latest`).

Prerequisites: Docker daemon running, clean working tree, and a valid login
to the registry (`docker login code.andydayton.com`). The script pre-flights
all three and aborts before pushing the git tag if any fail.
