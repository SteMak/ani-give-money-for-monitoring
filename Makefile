GO_BUILD_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
DOCKER_BUILD=$(shell pwd)/.docker_build
DOCKER_CMD_WORKER=$(DOCKER_BUILD)/worker

heroku: clean build
	heroku container:push web

clean:
	rm -rf $(DOCKER_BUILD)

build: $(DOCKER_CMD_WEB) $(DOCKER_CMD_WORKER)

$(DOCKER_BUILD):
	mkdir -p $(DOCKER_BUILD)

$(DOCKER_CMD_WORKER): $(DOCKER_BUILD)
	go build -v -o $(DOCKER_CMD_WORKER) cmd/worker/main.go
