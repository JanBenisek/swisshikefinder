# Specifies a parent image
# FROM golang:1.20.7-bullseye AS builder
FROM golang:1.20.7-bullseye
 
# Creates an app directory to hold your app’s source code
WORKDIR /app
 
# Copies everything from your root directory into /app
COPY /src .

# Installs Go dependencies
RUN go mod download

# Build the Go application into a binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /swiss-hiker-bin

# Use a lightweight Alpine image as the final base image
# FROM alpine:latest AS build-release-stage
# FROM gcr.io/distroless/base-debian11 AS build-release-stage

# WORKDIR /
# COPY --from=build-stage /swiss-hiker-bin /swiss-hiker-bin
# EXPOSE 3000

ENTRYPOINT ["/swiss-hiker-bin"]
# CMD ["swiss-hiker-bin"]