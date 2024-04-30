FROM golang:1.22.2 AS build-stage

# this will label the github package as public
LABEL org.opencontainers.image.source="https://github.com/janbenisek/swisshikefinder"

WORKDIR /app

# recommended to use ./ which forces current working directory (WORKDIR)
COPY /src ./
COPY /data/results ./data

# Download go modules
RUN go mod download

# Build
# with ./ it goes into workdir
# with / it does to the root of the container, so next to app
# need cgo for duckdb
RUN CGO_ENABLED=1 GOOS=linux go build -o ./swiss-hiker-bin


# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# seems like cgo=1 needs ubuntu to execute the binary
FROM ubuntu:latest AS build-release-stage

# plain ubuntu has no tls certificates (maybe build and move from base?)
RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /app

COPY --from=build-stage app/swiss-hiker-bin ./swiss-hiker-bin
COPY --from=build-stage app/data/ ./data/

EXPOSE 8080

# TODO: add no-root user and run with it
# USER nonroot

# because WORKDIR is /app, ./ works
ENTRYPOINT ["./swiss-hiker-bin"]