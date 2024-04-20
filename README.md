# Swiss Hike Finder

- run with  `go build -o ./swisshikerbin && ./swisshikerbin`
- alive at http://localhost:3000
- `docker build -t swiss-hike-finder:latest .`
- `docker run swiss-hike-finder:latest`
  - or detached `docker run -d -p 8080:8080 swiss-hike-finder:latest`
- `docker run -it --entrypoint /bin/bash swisshikefinder-web`
- test
- another test

## Readings

- https://freshman.tech/web-development-with-go/
- https://stackoverflow.com/questions/47270595/how-to-parse-json-string-to-struct
- https://docs.docker.com/language/golang/build-images/
- https://dev.to/willvelida/pushing-container-images-to-github-container-registry-with-github-actions-1m6b