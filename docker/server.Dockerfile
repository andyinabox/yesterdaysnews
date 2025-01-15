# We are using this distro to be compatable with
# jelastic cloud *shrug*
FROM jelastic/golang:1.23.4-almalinux-9 AS builder

# Set the working directory
WORKDIR /app

# Copy the Go source code
COPY . .

# download dependencies
RUN go mod download

# expose ports
EXPOSE 8080

# Build the Go binary
RUN go build -o /main .

FROM jelastic/almalinuxvps:9.3

WORKDIR /

COPY --from=builder /main /main

# Run the applicatio
CMD ["/main"]
