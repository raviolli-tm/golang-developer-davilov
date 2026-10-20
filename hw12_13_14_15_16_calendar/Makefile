BIN_DIR := ./bin
CALENDAR_BIN := $(BIN_DIR)/calendar
SCHEDULER_BIN := $(BIN_DIR)/calendar_scheduler
STORER_BIN := $(BIN_DIR)/calendar_storer

GIT_HASH := $(shell git log --format="%h" -n 1 2>/dev/null)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

.PHONY: build build-calendar build-scheduler build-storer run version test lint

build: build-calendar build-scheduler build-storer

build-calendar:
	@mkdir -p $(BIN_DIR)
	go build -v -o $(CALENDAR_BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

build-scheduler:
	@mkdir -p $(BIN_DIR)
	go build -v -o $(SCHEDULER_BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler

build-storer:
	@mkdir -p $(BIN_DIR)
	go build -v -o $(STORER_BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar_storer

run: build-calendar
	$(CALENDAR_BIN) -config ./configs/config.yaml

version: build-calendar
	$(CALENDAR_BIN) version

test:
	go test -race -count=1 ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.62.2

lint: install-lint-deps
	golangci-lint run ./...
