# TODO: use jelastic Debian image with backports to install yt-dlp

#
# build stage
#
# FROM golang:1.22.11 AS build

# # Set the working directory
# WORKDIR /app

# # Copy the Go source code
# COPY . .

# # download dependencies
# RUN go mod download

# # Build the Go binary
# RUN go build -o /main ./cmd/builder/main.go

# #
# # final stage
# #
# FROM jelastic/ubuntuvps:22.04

# WORKDIR /

# EXPOSE 80


# # install yt-dlp from binary
# RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
# RUN chmod a+rx /usr/local/bin/yt-dlp
# RUN /usr/local/bin/yt-dlp -U
# ENV YN_YT_DLP_PATH=/usr/local/bin/yt-dlp

# RUN apt update
# RUN apt install -y ffmpeg

# COPY --from=build /main /usr/local/bin/builder

# # Run the applicatio
# CMD ["/usr/local/bin/builder", "--output", "/dist", "-b", "setup"]


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

FROM builder-base

# Copy the source code
COPY . .

# download dependencies and build
RUN go mod download
RUN go build -o /usr/local/bin/builder ./cmd/builder/main.go

EXPOSE 80

CMD ["/usr/local/bin/builder", "--output", "/dist", "-v"]



# FROM golang:1.22.11 AS build

# # Set the working directory
# WORKDIR /app

# # Copy the Go source code
# COPY . .

# # download dependencies
# RUN go mod download

# # Build the Go binary
# RUN go build -o /main ./cmd/builder/main.go


# #
# # final stage
# #
# FROM jelastic/debianvps:12.9

# WORKDIR /

# # RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o ~/.local/bin/yt-dlp
# # RUN chmod a+rx ~/.local/bin/yt-dlp

# RUN printf "deb http://httpredir.debian.org/debian bookworm-backports main non-free\ndeb-src http://httpredir.debian.org/debian bookworm-backports main non-free" > /etc/apt/sources.list.d/backports.list
# # RUN add-apt-repository -y ppa:tomtomtom/yt-dlp
# RUN apt update
# RUN apt install -y yt-dlp -t bookworm-backports
# # RUN apt install -y yt-dlp
# RUN apt install -y ffmpeg
# # RUN yt-dlp -U

# COPY --from=build /main /usr/local/bin/builder

# EXPOSE 80

# # Run the applicatio
# CMD ["/usr/local/bin/builder", "--output", "/dist", "-b", "setup"]


#
# build stage
#
# FROM golang:1.22.11-alpine3.21 AS build
# FROM golang:1.22.11-alpine3.21

# # # Set the working directory
# WORKDIR /app

# # # Copy the Go source code
# COPY . .

# # # download dependencies
# RUN go mod download

# # # Build the Go binary
# RUN go build -o /main ./cmd/builder/main.go

# RUN apk update
# RUN apk add ffmpeg
# RUN apk add yt-dlp

# WORKDIR /

# EXPOSE 80

# # Run the applicatio
# CMD ["/main", "--output", "/dist"]

#
# final stage
#
# FROM alpine:3.21.2

# RUN apk update
# RUN apk add ffmpeg
# RUN apk add yt-dlp

# WORKDIR /

# COPY --from=build /main /main

# EXPOSE 80

# # Run the applicatio
# CMD ["/main", "--output", "/dist"]
