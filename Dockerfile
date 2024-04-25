FROM golang:1.20.7-bullseye AS build-stage

# this will label the github package as public
LABEL org.opencontainers.image.source="https://github.com/janbenisek/swisshikefinder"

WORKDIR /app

# recommended to use ./ which forces current working directory (WORKDIR)
COPY /src ./

# Download go modules
RUN go mod download

# Build
# with ./ it goes into workdir
# with / it does to the root of the container, so next to app
RUN CGO_ENABLED=0 GOOS=linux go build -o ./swiss-hiker-bin


# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM alpine:latest AS build-release-stage

WORKDIR /app

COPY --from=build-stage app/swiss-hiker-bin ./swiss-hiker-bin
# TODO: package static files into the binary
# COPY --from=build-stage app/assets/style.css ./assets/
# COPY --from=build-stage app/index.html ./

EXPOSE 8080

# TODO: add no-root user and run with it
# USER nonroot

# because WORKDIR is /app, ./ works
ENTRYPOINT ["./swiss-hiker-bin"]