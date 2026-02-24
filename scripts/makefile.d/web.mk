##@ Web
## FE sources
JS_SOURCES    := $(shell find assets/js             -type f -name '*.js'   -print -o \
                                                    -type f -name '*.jsx'  -print -o \
                                                    -type f -name '*.json' -print)
CSS_SOURCES   := $(shell find assets/css            -type f -name '*.css'  -print)
WEB_LIBS      := $(shell find assets/lib            -type f                -print)
HTML_SOURCES  := $(shell find assets/html-templates -type f -name '*.html' -print)
IMAGES        := $(shell find assets/images         -type f                -print)
WEB_META      := assets/manifest.json assets/favicon.ico

build/docker-image: $(JS_SOURCES) $(CSS_SOURCES) $(WEB_LIBS) $(HTML_SOURCES)

STATICS :=

## common static files
#STATICS += $(CSS_SOURCES:assets/%=build/static/%)
STATICS += $(WEB_LIBS:assets/%=build/static/%)
STATICS += $(IMAGES:assets/%=build/static/%)
#STATICS += $(WEB_META:assets/%=build/static/%)


build/static/%: assets/%
	@mkdir -p $(dir $@)
	cp $< $@

## esbuild

## js import helper
build/static/js/index.js: assets/js/index.js
assets/js/index.js: $(JS_SOURCES) scripts/tools/import-helper/*
	go run ./scripts/tools/import-helper > $@
OPTIONAL_CLEAN += assets/js/index.js

## esbuild
STATICS += build/static/js/index.js # entrypoints
build/static/js/%: export NODE_PATH=assets/js:assets/lib
build/static/js/%: $(JS_SOURCES) package.json package-lock.json | runtime/tools/esbuild
	@mkdir -p $(dir $@)
	esbuild $(@:build/static/%=assets/%) --outdir=$(dir $@) \
		--external:lit-html \
		--external:bm.js/bm.module.js \
		--minify \
		--bundle --sourcemap --format=esm \
		$(OPTIONAL_WEB_BUILD_ARGS)
	#--minify \
	#--external:/config.js


STATICS += build/static/css/page.css build/static/css/element.css
build/static/css/%.css: $(CSS_SOURCES) $(WEB_LIBS) | runtime/tools/esbuild
	@mkdir -p $(dir $@)
	esbuild $(@:build/static/%=assets/%) --outdir=$(dir $@) \
		--bundle --sourcemap \
		--minify \
		$(OPTIONAL_WEB_BUILD_ARGS)

gen: $(filter-out build/static/js/% build/static/css/%,$(STATICS))

.PHONY: build-web
build-web: $(STATICS) ## Build web-files. (bundle, minify, transpile, etc.)

build/$(APP_NAME): $(HTML_SOURCES) $(STATICS)

.PHONY: vet
vet: $(STATICS)

build/static:
	mkdir -p $@
	touch $@/.placeholder

## resolve depandancy
OPTIONAL_CLEAN += node_modules

build/$(APP_NAME): package-lock.json
package-lock.json: package.json | runtime/tools/npm
	@mkdir -p $(dir $@)
	npm install
	touch $@

build/docker-image: package.json

## Go esbuild (npm/npx 대신 사용)
runtime/tools/esbuild:
	GOBIN=$(shell pwd)/runtime/tools go install github.com/evanw/esbuild/cmd/esbuild@latest

build-tools: runtime/tools/esbuild runtime/tools/npm
runtime/tools/npm:
	@which $(notdir $@)
