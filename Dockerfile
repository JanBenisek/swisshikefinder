FROM golang:1.22.2 AS build-stage

# this will label the github package as public
LABEL org.opencontainers.image.source="https://github.com/janbenisek/swisshikefinder"

WORKDIR /app

# Copy go.mod and go.sum to cache dependencies
COPY /src/go.mod /src/go.sum ./

# Copy the entire project source, including the internal modules
COPY /src ./

# Download go modules
RUN go mod download

# copy the rest
# recommended to use ./ which forces current working directory (WORKDIR)
# COPY /src ./
COPY /data/db ./db
COPY /data/raw_data ./raw_data

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
# RUN apt-get update && apt-get install -y ca-certificates unzip wget

# TEMP TO TEST DUCKDB
# RUN mkdir -p /opt/duckdb && \
#     wget -O /opt/duckdb/duckdb_cli.zip https://github.com/duckdb/duckdb/releases/download/v0.10.2/duckdb_cli-linux-amd64.zip && \
#     unzip /opt/duckdb/duckdb_cli.zip -d /opt/duckdb && \
#     rm /opt/duckdb/duckdb_cli.zip

# Set the PATH environment variable to include DuckDB CLI
# ENV PATH="/opt/duckdb:${PATH}"

WORKDIR /app

COPY --from=build-stage app/swiss-hiker-bin ./swiss-hiker-bin
COPY --from=build-stage app/db/ ./db/
COPY --from=build-stage app/raw_data/ ./raw_data/

EXPOSE 8080

# TODO: add no-root user and run with it
# USER nonroot

# because WORKDIR is /app, ./ works
ENTRYPOINT ["./swiss-hiker-bin"]