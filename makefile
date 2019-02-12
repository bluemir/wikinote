IMPORT_PATH=github.com/bluemir/wikinote
BIN_NAME=$(notdir $(IMPORT_PATH))

DOCKER_IMAGE_NAME=bluemir/wikinote

default: $(BIN_NAME)

GIT_COMMIT_ID:=$(shell git rev-parse --short HEAD)
VERSION:=$(GIT_COMMIT_ID)-$(shell date +"%Y%m%d.%H%M%S")

# if gopath not set, make inside current dir
GO_SOURCES = $(shell find . -type f -name '*.go' -print)
JS_SOURCES = $(shell find app/js -type f -name '*.js' -print)
HTML_SOURCES = $(shell find app/html -type f -name '*.html' -print)
CSS_SOURCES = $(shell find app/less -type f -name "*.less" -print)
WEB_LIBS = $(shell find app/lib -type f -type f -print)

DISTS  = $(JS_SOURCES:app/js/%=dist/js/%)
DISTS += $(HTML_SOURCES:app/html/%=dist/html/%)
DISTS += dist/css/common.css
DISTS += $(WEB_LIBS:app/lib/%=dist/lib/%)

# Automatic runner
DIRS = $(shell find . -name dist -prune -o -name ".git" -prune -o -type d -print)

.sources:
	@echo $(DIRS) makefile \
		$(GO_SOURCES) \
		$(JS_SOURCES) \
		$(HTML_SOURCES) \
		$(CSS_SOURCES) \
		$(WEB_LIBS)| tr " " "\n"
run: $(BIN_NAME)
	./$(BIN_NAME) -D serve
auto-run:
	while true; do \
		make .sources | entr -rd make run ;  \
		echo "hit ^C again to quit" && sleep 1  \
	; done
reset:
	ps -e | grep make | grep -v grep | awk '{print $$1}' | xargs kill

## Binary build
$(BIN_NAME).bin: $(GO_SOURCES)
	go build -v \
		-ldflags "-X main.Version=$(VERSION)" \
		-o $(BIN_NAME).bin .
	@echo Build DONE

$(BIN_NAME): $(BIN_NAME).bin $(DISTS)
	cp $(BIN_NAME).bin $(BIN_NAME).tmp
	rice append -v --exec $(BIN_NAME).tmp \
		-i $(IMPORT_PATH)/pkgs/server  \
		-i $(IMPORT_PATH)/pkgs/renderer \
		-i $(IMPORT_PATH)/pkgs/config
	mv $(BIN_NAME).tmp $(BIN_NAME)
	@echo Embed resources DONE

## Web dist
dist/html/%.html: app/html/%.html
	@mkdir -p $(basename $@)
	cp $< $@
dist/css/common.css: $(CSS_SOURCES)
	lessc app/less/main.less $@
dist/%: app/%
	@mkdir -p $(basename $@)
	cp $< $@

tools:
	npm install -g less
	go get github.com/GeertJohan/go.rice/rice

clean:
	rm -rf dist/ vendor/ $(BIN_NAME) $(BIN_NAME).bin $(BIN_NAME).tmp
	go clean

docker-build: .docker-image
docker-push: .docker-image.pushed

.docker-image: Dockerfile makefile $(GO_SOURCES) $(DISTS)
	docker build -t $(DOCKER_IMAGE_NAME):$(VERSION) .
	echo "$(DOCKER_IMAGE_NAME):$(VERSION)" > .docker-image
.docker-image.pushed: .docker-image
	docker push $(shell cat .docker-image)
	echo $(shell cat .docker-image) > .docker-image.pushed

.PHONY: .sources run auto-run reset tools clean
