VERSION?=$(shell git describe --tags --dirty --always)
export VERSION

IMPORT_PATH=$(shell cat go.mod | head -n 1 | awk '{print $$2}')
APP_NAME=$(notdir $(IMPORT_PATH))

export GO111MODULE=on

export PATH:=${PATH}:./build/tools/

# go build args
OPTIONAL_BUILD_ARGS?=

default: build
	# `make help` for more information.

# sub-makefiles
# for build tools, docker build, deploy, static web files.
include scripts/makefile.d/*.mk

##@ General
clean: ## Clean up
	rm -rf build/ $(OPTIONAL_CLEAN)

build-tools: ## Install build tools
	# Build tool installed
tools: build-tools ## Install tools(include build tools)
	# Tool installed

help: ## Display this help
	# requirement
	#  - golang: 1.16.x
	#  - node  : 14.16.x
	#  - make  : 4.3 (*CAUTION* osx has lower verion of make)
	#
	@printf "# Usage:\n"
	@printf "#   make \033[36m<target>\033[0m\n"
	@awk 'BEGIN {FS = ":.*##";} /^[a-zA-Z_0-9-]+:.*?##/ { printf "#   \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "#\n# \033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: default build run test clean tools build-tools help

