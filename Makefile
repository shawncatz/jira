SHELL := /bin/bash

# The name of the executable (default is current directory name)
TARGET := $(shell echo $${PWD\#\#*/})
.DEFAULT_GOAL: $(TARGET)

# Github
GITHUB_USER := shawncatz
GITHUB_REPO := $(TARGET)

# These will be provided to the target
VERSION := $(shell cat version/VERSION)
TAG := v$(VERSION)
BUILD := `git rev-parse HEAD`
DATE := `date +%Y-%m-%dT%H:%M:%S%z`
PKG := github.com/$(GITHUB_USER)/$(GITHUB_REPO)/version

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X $(PKG).Build=$(BUILD) -X $(PKG).Version=$(VERSION) -X $(PKG).Date='$(DATE)'"

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all build clean install uninstall fmt simplify check run

all: check

$(TARGET): $(SRC)
	go build $(LDFLAGS) -o $(TARGET)

build: $(TARGET)
	@true

clean:
	@rm -f $(TARGET)

install:
	@go install $(LDFLAGS)

uninstall: clean
	@rm -f $$(which ${TARGET})

fmt:
	@gofmt -l -w $(SRC)

simplify:
	@gofmt -s -l -w $(SRC)

release-info:
	github-release info

release: clean tag-push build
	@echo creating release...
	@github-release release \
		--user $(GITHUB_USER) \
		--repo $(GITHUB_REPO) \
		--tag $(TAG) \
		--name "$(TARGET)"

	@echo uploading release...
	@github-release upload \
		--user $(GITHUB_USER) \
		--repo $(GITHUB_REPO) \
		--tag $(TAG) \
		--name "$(TARGET)" \
		--file $(TARGET)

tag:
	@echo creating tag...
	@-git tag $(TAG)

tag-push: tag
	@git push --tags

check:
	@test -z $(shell gofmt -l main.go | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d}; done
	@go tool vet ${SRC}

run: install
	@$(TARGET)
