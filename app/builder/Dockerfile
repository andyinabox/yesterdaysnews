FROM alpine:3.21.2

RUN apk upgrade
RUN apk add ffmpeg
RUN apk add yt-dlp

WORKDIR /

COPY bin/builder-linux-amd64 /main

# Run the applicatio
ENTRYPOINT ["/main"]
