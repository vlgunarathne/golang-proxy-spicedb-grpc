# syntax=docker/dockerfile:1

FROM golang:1.19 AS builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY . .
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
# COPY . .

# Build
RUN CGO_ENABLED=0 go build -o golang-proxy-spicedb

FROM golang:1.19

WORKDIR /app

COPY --from=builder /app/golang-proxy-spicedb .
# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 9090

# Run
CMD ["./golang-proxy-spicedb"]