FROM golang:1.22.4 AS build

# Set the working directory
WORKDIR /app

# Copy the Go source code
COPY . .

# download dependencies
RUN go mod download

# expose ports
EXPOSE 8080

# Build the Go binary
RUN GOOS=linux GOARCH=amd64 go build -o /main .

FROM jelastic/almalinuxvps:9.3

WORKDIR /

COPY --from=build /main /main

# Run the applicatio
CMD ["/main"]
