# syntax=docker/dockerfile:1

FROM golang:latest AS builder

ARG IMAGE_TAG

# Set destination for COPY
WORKDIR /app

# Copy .go files
COPY *.go ./

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Build
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags "-X 'main.IMAGE_TAG=$IMAGE_TAG'" -o /app/main .

# Stage 2
FROM alpine:latest

# Install nano
RUN apk update
RUN apk add nano

# Set the working directory
WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /app/main .

# Enable ip routing
RUN echo 1 > /proc/sys/net/ipv4/ip_forward

# Run
CMD ["/app/main"]