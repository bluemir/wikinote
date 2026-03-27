##@ Web

## FE sources
JS_SOURCES    := $(shell find assets/src/js         -type f -name '*.js' ! -name 'index.js' -print -o \
                                                    -type f -name '*.jsx'  -print -o \
                                                    -type f -name '*.json' -print)
CSS_SOURCES   := $(shell find assets/src/css        -type f -name '*.css'  -print)
HTML_SOURCES  := $(shell find assets/html-templates -type f -name '*.html' -print)

build/docker-image: $(JS_SOURCES) $(CSS_SOURCES) $(HTML_SOURCES)

## js import helper
OPTIONAL_CLEAN += assets/src/js/index.js
assets/src/js/index.js: $(JS_SOURCES) scripts/tools/import-helper/*
	go run ./scripts/tools/import-helper -o $@ assets/src/js

## dev build
build/$(APP_NAME): assets/src/js/index.js

## prod build: esbuild runs via go:generate in assets/prod.go
assets/dist/js/index.js: $(JS_SOURCES) $(CSS_SOURCES) assets/src/js/index.js
	go generate -tags prod -x ./...

gen: assets/dist/js/index.js

## resolve dependency
OPTIONAL_CLEAN += node_modules

build/$(APP_NAME):            package-lock.json
build/$(APP_NAME)-$(VERSION): package-lock.json

package-lock.json: package.json | runtime/tools/npm
	@mkdir -p $(dir $@)
	npm install
	touch $@

build/docker-image: package.json

build-tools: runtime/tools/npm
runtime/tools/npm:
	@which $(notdir $@)

## prod build uses dist folder
OPTIONAL_CLEAN += assets/dist
