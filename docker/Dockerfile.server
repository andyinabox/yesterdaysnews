# We are using this distro to be compatable with
# jelastic cloud *shrug*
FROM jelastic/golang:1.22.1-almalinux-9

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

# Run the applicatio
CMD ["/main"]
