# Swiss Hike Finder
Playground to learn Go and find awesome hikes.

## Handy commands

```shell
go build -o ./swisshikerbin && ./swisshikerbin
docker build -t swiss-hike-finder:latest .
docker compose up web
docker stop $(docker ps -a -q)
docker run -it --entrypoint /bin/bash swiss-hike-finder:latest
docker run -p 8080:8080 -e "HIKE_API_KEY=XXX" swiss-hike-finder:latest
```

## Readings

- https://freshman.tech/web-development-with-go/
- https://stackoverflow.com/questions/47270595/how-to-parse-json-string-to-struct
- https://docs.docker.com/language/golang/build-images/
- https://dev.to/willvelida/pushing-container-images-to-github-container-registry-with-github-actions-1m6b