VERSION?=$(shell git describe --tags --dirty --always)
export VERSION

IMPORT_PATH=$(shell cat go.mod | head -n 1 | awk '{print $$2}')
APP_NAME=$(notdir $(IMPORT_PATH))

export GO111MODULE=on

## Go Sources
GO_SOURCES = $(shell find . -name "vendor"  -prune -o \
                            -type f -name "*.go" -print)

## FE sources
JS_SOURCES    := $(shell find static/js             -type f -name '*.js'   -print)
CSS_SOURCES   := $(shell find static/css            -type f -name '*.css'  -print)
WEB_LIBS      := $(shell find static/lib            -type f                -print)
HTML_SOURCES  := $(shell find static/html-templates -type f -name '*.html' -print)

STATICS :=
STATICS += $(JS_SOURCES:%=build/%)
STATICS += $(CSS_SOURCES:%=build/%)
STATICS += $(WEB_LIBS:%=build/%)

## see Makefile.d/nodejs.mk for using rollup, less or other tools

default: build

# sub-makefiles
# for build tools, docker build, deploy
include makefile.d/*

## Static files
build/static/%: static/%
	@mkdir -p $(dir $@)
	cp $< $@

build: build/$(APP_NAME)

build/$(APP_NAME).unpacked: $(GO_SOURCES) $(MAKEFILE_LIST)
	@$(MAKE) build/tools/go
	@mkdir -p build
	go build -v \
		-trimpath \
		-ldflags "\
			-X main.AppName=$(APP_NAME) \
			-X main.Version=$(VERSION)  \
		" \
		$(OPTIONAL_BUILD_ARGS) \
		-o $@ main.go

build/$(APP_NAME): build/$(APP_NAME).unpacked $(HTML_SOURCES) $(STATICS) $(MAKEFILE_LIST)
	$(MAKE) build/tools/rice
	@mkdir -p $(dir $<)
	cp $< $@.tmp
	rice append -v \
		-i $(IMPORT_PATH)/pkg/static \
		--exec $@.tmp
	mv $@.tmp $@

clean:
	rm -rf build/ $(OPTIONAL_CLEAN_DIR)

run: build/$(APP_NAME)
	$< -vvvv server

auto-run:
	while true; do \
		$(MAKE) .watched_sources | \
		entr -rd $(MAKE) test run ;  \
		echo "hit ^C again to quit" && sleep 1  \
	; done

reset:
	ps -e | grep make | grep -v grep | awk '{print $$1}' | xargs kill

## watched_sources
.watched_sources: \
	$(MAKEFILE_LIST) \
	go.mod go.sum \
	$(GO_SOURCES) \
	$(JS_SOURCES) \
	$(CSS_SOURCES) \
	$(WEB_LIBS) \
	$(HTML_SOURCES)
	@echo $^ | tr " " "\n"

test:
	go test -v ./pkg/...

.PHONY: build clean run auto-run reset .watched_sources test
