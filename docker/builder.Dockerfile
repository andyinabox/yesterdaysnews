#
# build stage
#
# FROM golang:1.22.11-alpine3.21 AS build
FROM golang:1.22.11-alpine3.21

# # Set the working directory
WORKDIR /app

# # Copy the Go source code
COPY . .

# # download dependencies
RUN go mod download

# # Build the Go binary
RUN go build -o /main ./cmd/build/main.go

RUN apk update
RUN apk add ffmpeg
RUN apk add yt-dlp

WORKDIR /

EXPOSE 80

# Run the applicatio
CMD ["/main", "--output", "/dist"]

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
