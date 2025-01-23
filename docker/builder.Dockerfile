#
# build stage
#
FROM golang:1.22.11-alpine3.21 AS build

# # Set the working directory
WORKDIR /app

# # Copy the Go source code
COPY . .

# # download dependencies
RUN go mod download

# # Build the Go binary
RUN go build -o /main ./cmd/builder/main.go

#
# base stage (installs dependencies)
#
FROM alpine:3.21.2 AS builder-base

RUN apk update
RUN apk add ffmpeg
RUN apk add yt-dlp

#
# final stage
#
FROM builder-base

WORKDIR /

COPY --from=build /main /main

EXPOSE 80

# Run the applicatio
ENTRYPOINT ["/main"]
