#
# build stage
#
FROM golang:1.22.4 AS build

# Set the working directory
WORKDIR /app

# Copy the Go source code
COPY . .

# download dependencies
RUN go mod download

# Build the Go binary
RUN GOOS=linux GOARCH=amd64 go build -o /main .

#
# final stage
#
FROM jelastic/almalinuxvps:9.3

WORKDIR /

COPY --from=build /main /main

EXPOSE 8080

# Run the applicatio
CMD ["/main"]
