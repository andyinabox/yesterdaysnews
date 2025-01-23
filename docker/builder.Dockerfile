#
# stage 1
#
FROM jelastic/golang:1.22.11-almalinux-9 AS builder-base

# install ffmpeg
RUN dnf install -y epel-release
RUN dnf config-manager --set-enabled crb
RUN dnf install --nogpgcheck https://mirrors.rpmfusion.org/free/el/rpmfusion-free-release-$(rpm -E %rhel).noarch.rpm -y
RUN dnf install --nogpgcheck https://mirrors.rpmfusion.org/nonfree/el/rpmfusion-nonfree-release-$(rpm -E %rhel).noarch.rpm -y
RUN dnf install -y ffmpeg

# install yt-dlp from binary
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
RUN chmod a+rx /usr/local/bin/yt-dlp
RUN /usr/local/bin/yt-dlp -U
ENV YN_YT_DLP_PATH=/usr/local/bin/yt-dlp


#
# stage 2
#
FROM builder-base

# Copy the source code
COPY . .

# download dependencies and build
RUN go mod download
RUN go build -o /usr/local/bin/builder ./cmd/builder/main.go

EXPOSE 80

CMD ["/usr/local/bin/builder", "--output", "/dist", "-v"]

