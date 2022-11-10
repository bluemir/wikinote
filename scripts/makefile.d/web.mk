##@ Web
## FE sources
JS_SOURCES    := $(shell find web/js             -type f -name '*.js'   -print -o \
                                                 -type f -name '*.jsx'  -print -o \
                                                 -type f -name '*.json' -print)
CSS_SOURCES   := $(shell find web/css            -type f -name '*.css'  -print)
WEB_LIBS      := $(shell find web/lib            -type f                -print)
HTML_SOURCES  := $(shell find web/html-templates -type f -name '*.html' -print)
IMAGES        := $(shell find web/images         -type f                -print)
WEB_META      := web/manifest.json web/favicon.ico

.watched_sources: $(JS_SOURCES) $(CSS_SOURCES) $(WEB_LIBS) $(HTML_SOURCES)
build/docker-image: $(JS_SOURCES) $(CSS_SOURCES) $(WEB_LIBS) $(HTML_SOURCES)

STATICS :=

## common static files
STATICS += $(CSS_SOURCES:web/%=build/static/%)
STATICS += $(WEB_LIBS:web/%=build/static/%)
STATICS += $(IMAGES:web/%=build/static/%)
#STATICS += $(WEB_META:web/%=build/static/%)


build/static/%: web/%
	@mkdir -p $(dir $@)
	cp $< $@

## esbuild
STATICS += build/static/js/v1/index.js # entrypoint
build/static/js/%: export NODE_PATH=web/js:web/lib
build/static/js/%: $(JS_SOURCES) build/yarn-updated
	@$(MAKE) build/tools/npx
	@mkdir -p $(dir $@)
	npx esbuild $(@:build/static/%=web/%) --outdir=$(dir $@) \
		--bundle --sourcemap --format=esm $(OPTIONAL_WEB_BUILD_ARGS)
	#--external:/config.js \
	#--minify \

.PHONY: build-web
build-web: $(STATICS) ## Build web-files. (bundle, minify, transpile, etc.)

build/$(APP_NAME): $(HTML_SOURCES) $(STATICS)

## resolve depandancy
OPTIONAL_CLEAN += node_modules

build/$(APP_NAME): build/yarn-updated
build/yarn-updated: package.json
	@$(MAKE) build/tools/yarn
	@mkdir -p $(dir $@)
	yarn install
	touch $@

.watched_sources: package.json
build/docker-image: package.json

build-tools: build/tools/npm build/tools/yarn build/tools/npx
build/tools/npm:
	@which $(notdir $@)
build/tools/npx:
	@which $(notdir $@)
build/tools/yarn: build/tools/npm
	@which $(notdir $@) || (npm install -g yarn)
