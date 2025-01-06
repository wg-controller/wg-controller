# syntax=docker/dockerfile:1

FROM golang:latest AS builder
ARG IMAGE_TAG

# Set destination for COPY
WORKDIR /app

# Copy .go files
COPY *.go ./
COPY db/ db/
COPY types/ types/

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Set env vars
ENV CGO_ENABLED=1
ENV GOOS=linux

# Build
RUN go build -ldflags "-X 'main.IMAGE_TAG=$IMAGE_TAG' -linkmode external -extldflags '-static'" -o /app/main .


# Final stage
FROM alpine:3.14.2

# Install networking packages
RUN apk add --no-cache \
    dpkg \
    dumb-init \
    iptables \
    iptables-legacy \
    wireguard-tools

# Copy binaries from build stage
COPY --from=builder /app/main /app/main

# Run
CMD ["/app/main"]