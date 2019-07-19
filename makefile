IMPORT_PATH=github.com/bluemir/wikinote
BIN_NAME=$(notdir $(IMPORT_PATH))

DOCKER_IMAGE_NAME=bluemir/wikinote

default: build/$(BIN_NAME)

VERSION?=$(shell git describe --long --tags --dirty --always)

GO_SOURCES   = $(shell find .        -type f -name '*.go'   -print)
JS_SOURCES   = $(shell find app/js   -type f -name '*.js'   -print)
HTML_SOURCES = $(shell find app/html -type f -name '*.html' -print)
CSS_SOURCES  = $(shell find app/less -type f -name "*.less" -print)
WEB_LIBS     = $(shell find app/lib  -type f                -print)

DISTS  = $(JS_SOURCES:app/js/%=build/dist/js/%)
DISTS += $(HTML_SOURCES:app/html/%=build/dist/html/%)
DISTS += build/dist/css/common.css
DISTS += $(WEB_LIBS:app/lib/%=build/dist/lib/%)

DIRS = $(shell find . \
                    -name build -prune -o \
                    -name ".git" -prune -o \
                    -type d \
                    -print)
.sources:
	@echo $(DIRS) makefile \
	      $(GO_SOURCES) \
	      $(JS_SOURCES) \
	      $(HTML_SOURCES) \
	      $(CSS_SOURCES) \
	      $(WEB_LIBS)| tr " " "\n"

# BUILD

## Binary build
build/$(BIN_NAME).bin: $(GO_SOURCES) makefile
	go build -v \
		-ldflags "-X main.Version=$(VERSION)" \
		-o $@ .
	@echo Build DONE

## Web dist
build/dist/html/%.html: app/html/%.html
	@mkdir -p $(dir $@)
	cp $< $@
build/dist/css/common.css: $(CSS_SOURCES)
	lessc app/less/main.less $@
build/dist/%: app/%
	@mkdir -p $(dir $@)
	cp $< $@

## resource embed
build/$(BIN_NAME): build/$(BIN_NAME).bin $(DISTS)
	cp $< $@.tmp
	rice append -v \
		-i $(IMPORT_PATH)/pkgs/dist \
		--exec $@.tmp
	mv $@.tmp $@
	@echo Embed resources DONE

test:
	go test -v ./...

run: export LOG_LEVEL=trace
run: build/$(BIN_NAME)
	build/$(BIN_NAME) -D serve
auto-run:
	while true; do \
		make .sources | entr -rd make test run ;  \
		echo "hit ^C again to quit" && sleep 1  \
	; done
reset:
	ps -f -C make | grep "test run" | awk '{print $$2}' | xargs kill

docker-build: build/.docker-image
build/.docker-image: Dockerfile makefile $(GO_SOURCES) $(DISTS)
	docker build \
		--build-arg VERSION=$(VERSION) \
		-t $(DOCKER_IMAGE_NAME):$(VERSION) .
	echo "$(DOCKER_IMAGE_NAME):$(VERSION)" > $@


docker-push: build/.docker-image.pushed
build/.docker-image.pushed: .docker-image
	docker push $(shell cat .docker-image)
	echo $(shell cat .docker-image) > $@

docker-run: build/.docker-image
	docker run --rm -it \
		-p 4000:4000 \
		-v ~/wiki:/wiki \
		-e LOG_LEVEL=trace \
		$(shell cat $<) \
		serve \
		--wiki-path /wiki \
		--config /wiki/.app/config.yaml \
		--bind :4000

tools:
	npm install -g less
	go get github.com/GeertJohan/go.rice/rice

clean:
	rm -rf bulid/ vendor/
	go clean

.PHONY: .sources run auto-run reset docker-build docker-push docker-run tools clean
