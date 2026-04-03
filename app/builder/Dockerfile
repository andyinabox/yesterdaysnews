FROM alpine:3.22.3

ARG TARGETARCH

RUN apk upgrade
RUN apk add ffmpeg nodejs

RUN if [ "$TARGETARCH" = "arm64" ]; then \
      wget -O /usr/local/bin/yt-dlp https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_musllinux_aarch64; \
    else \
      wget -O /usr/local/bin/yt-dlp https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_musllinux; \
    fi && chmod +x /usr/local/bin/yt-dlp
RUN echo "--js-runtimes node" > /etc/yt-dlp.conf

WORKDIR /

COPY bin/builder-linux-${TARGETARCH} /main

ENTRYPOINT ["/main"]
