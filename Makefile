.PHONY: all build lint coverage

all: build

# Project name
PROJECT_NAME := cineguard
GO_MODULE_NAME := "github.com/kaiokei/$(PROJECT_NAME)"

# Useful variables for build metadata
VERSION ?= $(shell git describe --tags --always --dirty)
# equivalent command to test git dirty status in bash terminal: [[ -n "$(git status --porcelain)" ]] && echo "true" || echo "false"
IS_GIT_DIRTY := $(shell [ -n "$$(git status --porcelain)" ] && echo "true" || echo "false")
COMMIT_LONG ?= $(shell git rev-parse HEAD)
COMMIT_SHORT ?= $(shell git rev-parse --short=8 HEAD)
COMMIT_TIMESTAMP := $(shell git show -s --format=%cI HEAD)
GO_VERSION ?= $(shell go version)
BUILD_PLATFORM  ?= $(shell uname -m)
BUILD_DATE ?= $(shell date -u --iso-8601=seconds)

LDFLAGS = "-X '$(GO_MODULE_NAME)/pkg/version.RawGitDescribe=$(VERSION)' \
	-X '$(GO_MODULE_NAME)/pkg/version.GitCommitIdLong=$(COMMIT_LONG)' \
	-X '$(GO_MODULE_NAME)/pkg/version.GitCommitIdShort=$(COMMIT_SHORT)' \
	-X '$(GO_MODULE_NAME)/pkg/version.GoVersion=$(GO_VERSION)' \
	-X '$(GO_MODULE_NAME)/pkg/version.BuildPlatform=$(BUILD_PLATFORM)' \
	-X '$(GO_MODULE_NAME)/pkg/version.BuildDate=$(BUILD_DATE)' \
	-X '$(GO_MODULE_NAME)/pkg/version.GitCommitTimestamp=$(COMMIT_TIMESTAMP)' \
	-X '$(GO_MODULE_NAME)/pkg/version.GitDirtyStr=$(IS_GIT_DIRTY)'"

GO_LDFLAGS = -ldflags=$(LDFLAGS)
BINARY_NAME = $(PROJECT_NAME)

## Pipeline

# Go parameters
CGO_ENABLED := "1"

lint:
		@golangci-lint run
coverage:
		mkdir -p build
		go test -race -v -coverprofile build/coverage.out ./pkg/...
		go tool cover -html=build/coverage.out -o build/coverage.html

## Dev
build:
		@go version
		@go build $(GO_LDFLAGS) -o $(PROJECT_NAME) cmd/main.go
build-debug:
		@go version
		@go build -gcflags="all=-N -l" $(GO_LDFLAGS) -o $(PROJECT_NAME) cmd/main.go
		$(info use cmd : dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec $(PROJECT_NAME))
		$(info will listen to port 2345)

run-test:
		@go run cmd/main.go test

## Docs
doc:
		@go run $(GO_LDFLAGS) cmd/main.go docs --output-dir docs/cli-user-interface/markdown/
		@go run $(GO_LDFLAGS) cmd/main.go docs --output-dir docs/cli-user-interface/txt/ --format cli-table-pretty

## Deploy

release: 
		@echo "Makefile: Running goreleaser release --clean fro project $(PROJECT_NAME)"
		LDFLAGS=$(LDFLAGS) goreleaser release --clean --skip sign,validate,ko
get-ldflags:
		@echo $(LDFLAGS)
