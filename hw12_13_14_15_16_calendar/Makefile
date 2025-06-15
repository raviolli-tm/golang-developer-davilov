BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"
APP_NAME="calendar-image"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

generate:
	docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli generate \
      -i /local/api/swagger.yaml \
      -g go-server \
      -o /local/api \
      -c /local/api/s.json \
		  --additional-properties=packageName=api,outputAsLibrary=true \




build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	@echo "Stop and delete old docker container(if exists)..."
	-docker stop $(APP_NAME) 2>/dev/null || true
	-docker rm $(APP_NAME) 2>/dev/null || true

	@echo "Run the new docker container..."
	docker run -d --name $(APP_NAME) -p 8080:8080 $(DOCKER_IMG)


version: build
	$(BIN) version

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint
