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

# Build wg-controller
RUN go build -ldflags "-X 'main.IMAGE_TAG=$IMAGE_TAG' -linkmode external -extldflags '-static'" -o /app/main .

# Download wireguard-go
RUN git clone https://git.zx2c4.com/wireguard-go
WORKDIR /app/wireguard-go

# Build wireguard-go
RUN go build -o /app/wireguard-go/wireguard-go .

# UI build stage
FROM node:lts-alpine AS ui-build

WORKDIR /app

# Copy package.json and package-lock.json
COPY ./wg-controller-ui/package*.json ./
COPY ./wg-controller-ui/ ./

# install project dependencies
RUN npm install

# build app
RUN npm run build

# Final stage
FROM alpine:3.14.2

# Install required packages
RUN apk add --no-cache bash libmnl iptables openresolv iproute2 libc6-compat

# Copy binaries from build stage
COPY --from=builder /app/main /app/main
COPY --from=builder /app/wireguard-go/wireguard-go /usr/bin/wireguard-go

# Install dnsmasq
RUN apk add --no-cache dnsmasq

# Packages to help with debugging
RUN apk add --no-cache -U wireguard-tools

# Copy UI build
COPY --from=ui-build /app/dist /var/www

# Run
ENTRYPOINT ["/app/main"]