## FE sources
JS_SOURCES    := $(shell find web/js             -type f -name '*.js'   -print)
CSS_SOURCES   := $(shell find web/css            -type f -name '*.css'  -print)
WEB_LIBS      := $(shell find web/lib            -type f                -print)
HTML_SOURCES  := $(shell find web/html-templates -type f -name '*.html' -print)

.watched_sources: $(JS_SOURCES) $(CSS_SOURCES) $(WEB_LIBS) $(HTML_SOURCES)
build/docker-image: $(JS_SOURCES) $(CSS_SOURCES) $(WEB_LIBS) $(HTML_SOURCES)


STATICS :=

## common static files
STATICS += $(CSS_SOURCES:web/%=build/static/%)
STATICS += $(WEB_LIBS:web/%=build/static/%)
build/static/%: web/%
	@mkdir -p $(dir $@)
	cp $< $@

## rollup & js
STATICS += build/static/js/v1/index.js                   # entrypoint
build/static/js/%: $(JS_SOURCES) build/yarn-updated
	@$(MAKE) build/tools/npx
	@mkdir -p $(dir $@)
	npx rollup $(@:build/static/%=web/%) --file $@ --format es -m -p '@rollup/plugin-node-resolve'

## less
## yarn add --dev less
#LESS_SOURCES  = $(shell find web/less           -type f -name '*.less' -print)
#STATICS := $(filter-out build/static/css/%,$(STATICS)) # remove default css files
#STATICS += $(LESS_SOURCES:web/less/%=build/static/css/%)
#build/static/css/%: web/less/% build/yarn-updated
#	@$(MAKE) build/tools/npx
#	@mkdir -p $(dir $@)
#	npx lessc $< $@
#.watched_sources: $(LESS_SOURCES)


build/$(APP_NAME): $(HTML_SOURCES) $(STATICS)

## resolve depandancy
OPTIONAL_CLEAN_DIR += node_modules

build/$(APP_NAME): build/yarn-updated
build/yarn-updated: package.json
	@$(MAKE) build/tools/yarn
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
