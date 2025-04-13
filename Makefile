BINARY_NAME=grafana-annotator
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS=-ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

.PHONY: build
build:
	go build ${LDFLAGS} -o bin/${BINARY_NAME} cmd/grafana-annotator/main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	rm -f bin/${BINARY_NAME}

.PHONY: install
install:
	go install ${LDFLAGS} ./cmd/grafana-annotator

.PHONY: lint
lint:
	golangci-lint run

.PHONY: all
all: clean build test

.DEFAULT_GOAL := build
