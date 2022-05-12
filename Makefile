DOCKER_ACCOUNT = vsvegner

PWD := $(shell pwd)
PROJECTNAME = $(shell basename $(PWD))
PROGRAM_NAME = $(shell basename $(PWD))

VERSION=$(shell git describe --tags)
COMMIT=$(shell git rev-parse --short HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
TAG=$(shell git describe --tags |cut -d- -f1)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

PLATFORMS=linux windows
# PLATFORMS=darwin linux windows
# ARCHITECTURES=386 amd64 ppc64 arm arm64
ARCHITECTURES=386 amd64 arm arm64

LDFLAGS = -ldflags "-s -w -X=main.Version=${VERSION} -X=main.Build=${COMMIT} -X main.gitTag=${TAG} -X main.gitCommit=${COMMIT} -X main.gitBranch=${BRANCH} -X main.buildTime=${BUILD_TIME}"

# Check for required command tools to build or stop immediately
EXECUTABLES = git go find pwd basename
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH)))

.PHONY: help clean dep build install uninstall pack

.DEFAULT_GOAL := help

help: ## Display this help screen.
	@echo "Makefile available targets:"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  * \033[36m%-15s\033[0m %s\n", $$1, $$2}'

clean: ## Clean bin directory.
	rm -f ./bin/*

dep: ## Download the dependencies.
	go mod tidy
	go mod download
#	go mod vendor

build: ## Build program executable for linux platform.
	mkdir -p ./bin
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o bin/${PROGRAM_NAME}_$(VERSION)_linux_$(COMMIT)_amd64 .
	sudo chmod +x bin/${PROGRAM_NAME}_$(VERSION)_linux_$(COMMIT)_amd64

pack: ## Packing all executable files in ./bin using UPX 
	upx ./bin/${PROGRAM_NAME}*

install: ## Install program executable into /usr/bin directory.
	mkdir -p /usr/bin/${PROGRAM_NAME}
	install -pm 755 bin/${PROGRAM_NAME} /usr/bin/${PROGRAM_NAME}/${PROGRAM_NAME}
	cp config.yaml.example /usr/bin/config.yaml

uninstall: ## Uninstall program executable from /usr/bin directory.
	rm -rf /usr/bin/${PROGRAM_NAME}
