include build.env
include .env

.DEFAULT_GOAL := build

build:
	DOCKER_BUILDKIT=1 docker build -t $(PROJECT_NAME):latest -f Dockerfile .

gobuild:
	go build -o ./swisshikerbin && ./swisshikerbin

compose:
	docker compose up -d && docker ps -a && docker logs $(PROJECT_NAME)-$(PROJECT_NAME)-1 -f

stop:
	docker stop $$(docker ps -a -q) && docker rm $$(docker ps -a -q)

rebuild: stop build compose

clean:
	docker system prune -af 
	docker system prune --volumes -f

run_bash:
	docker run --rm -it --entrypoint /bin/sh $(PROJECT_NAME):latest

run_port:
	docker run -p $(PORT):$(PORT) -e "HIKE_API_KEY=$(HIKE_API_KEY)" --rm -it $(PROJECT_NAME):latest

run_package:
	docker pull ghcr.io/janbenisek/swisshikefinder:latest
	docker run -p $(PORT):$(PORT) -e "HIKE_API_KEY=$(HIKE_API_KEY)" --rm -it ghcr.io/janbenisek/swisshikefinder:latest



.PHONY: build gobuild compose stop rebuild clean run_bash run_port run_package