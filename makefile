VERSION?=$(shell git describe --tags --dirty --always)
GIT_COMMIT_ID?=$(shell git rev-parse --short HEAD)
export VERSION
export GIT_COMMIT_ID

IMPORT_PATH=$(shell cat go.mod | head -n 1 | awk '{print $$2}')
BIN_NAME=$(notdir $(IMPORT_PATH))

export GO111MODULE=on
export GIT_TERMINAL_PROMPT=1

DOCKER_IMAGE_NAME=bluemir/$(BIN_NAME)

## Go Sources
GO_SOURCES = $(shell find . -name "vendor"  -prune -o \
                            -type f -name "*.go" -print)

## FE sources
JS_SOURCES    = $(shell find app/js       -type f -name '*.js'   -print)
HTML_SOURCES  = $(shell find app/html     -type f -name '*.html' -print)
CSS_SOURCES   = $(shell find app/css      -type f -name '*.css'  -print)
WEB_LIBS      = $(shell find app/lib      -type f                -print)

DISTS =
DISTS += $(JS_SOURCES:app/js/%=build/dist/js/%)
DISTS += $(CSS_SOURCES:app/css/%=build/dist/css/%)
DISTS += $(WEB_LIBS:app/lib/%=build/dist/lib/%)

HTML_TEMPLATE = $(HTML_SOURCES:app/html/%=build/template/%)

default: build

## Web dist
build/dist/$(GIT_COMMIT_ID)/%: app/%
build/dist/%: app/%
	@mkdir -p $(dir $@)
	cp $< $@
#dist/css/%.css: $(CSS_SOURCES)
#	lessc app/less/entry/$*.less $@
## HTML template
build/template/%: app/html/%
	@mkdir -p $(dir $@)
	cat $< | GIT_COMMIT_ID=$(GIT_COMMIT_ID) envsubst '$$GIT_COMMIT_ID' > $@

build: build/$(BIN_NAME)

build/$(BIN_NAME).unpacked: $(GO_SOURCES) makefile
	@mkdir -p build
	go build -v \
		-ldflags "-X main.VERSION=$(VERSION) -X main.GitCommitId=$(GIT_COMMIT_ID)" \
		$(OPTIONAL_BUILD_ARGS) \
		-o $@ main.go
build/$(BIN_NAME): build/$(BIN_NAME).unpacked $(HTML_TEMPLATE) $(DISTS)
	@mkdir -p build
	cp $< $@.tmp
	rice append -v \
		-i $(IMPORT_PATH)/pkg/dist \
		--exec $@.tmp
	mv build/$(BIN_NAME).tmp $@

docker: build/.docker-image

build/.docker-image: build/Dockerfile $(GO_SOURCES) $(HTML_TEMPLATE) $(DISTS)
	docker build \
		--build-arg VERSION=$(VERSION) \
		-t $(DOCKER_IMAGE_NAME):$(VERSION) \
		-f $< .
	echo $(DOCKER_IMAGE_NAME):$(VERSION) > $@

build/Dockerfile: export BIN_NAME:=$(BIN_NAME)
build/Dockerfile: Dockerfile.template
	@mkdir -p build
	cat $< | envsubst '$$BIN_NAME' > $@


push: build/.docker-image.pushed

build/.docker-image.pushed: build/.docker-image
	docker push $(shell cat build/.docker-image)
	echo $(shell cat build/.docker-image) > $@

clean:
	rm -rf build/

run: export LOG_LEVEL=TRACE
run: build/$(BIN_NAME)
	$< --admin-user root=1234 -w ./build/test-data -c ./build/test-data/.app/config.yaml serve

auto-run:
	while true; do \
		$(MAKE) .sources | \
		entr -rd $(MAKE) run ;  \
		echo "hit ^C again to quit" && sleep 1  \
	; done

.sources:
	@echo \
	makefile \
	$(GO_SOURCES) \
	$(JS_SOURCES) \
	$(HTML_SOURCES) \
	$(CSS_SOURCES) \
	$(WEB_LIBS) \
	$(HTML_TEMPLATE) \
	| tr " " "\n"

test:
	go test -v ./pkg/...

.PHONY: build docker push clean run auto-run .sources test default
