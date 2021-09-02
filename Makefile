VERSION?=$(shell git describe --tags --dirty --always)
export VERSION

IMPORT_PATH=$(shell cat go.mod | head -n 1 | awk '{print $$2}')
APP_NAME=$(notdir $(IMPORT_PATH))

export GO111MODULE=on

# go build args
OPTIONAL_BUILD_ARGS :=

default: build

# sub-makefiles
# for build tools, docker build, deploy, static web files.
include scripts/makefile.d/*

build: build/$(APP_NAME)

run: build/$(APP_NAME)
	$< -vvv server --admin-user root=1234 --wiki-path=runtime --config=runtime/.app/config.yaml

test:
	go test ./...

clean:
	rm -rf build/ $(OPTIONAL_CLEAN_DIR)

tools: build-tools
	@echo "--- done ---"


.PHONY: default build run test clean tools build-tools

