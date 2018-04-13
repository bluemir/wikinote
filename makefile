IMPORT_PATH=github.com/bluemir/wikinote
BIN_NAME=$(notdir $(IMPORT_PATH))

default: $(BIN_NAME)

GIT_COMMIT_ID:=$(shell git rev-parse --short HEAD)
VERSION:=$(GIT_COMMIT_ID)-$(shell date +"%Y%m%d.%H%M%S")

# if gopath not set, make inside current dir
ifeq ($(GOPATH),)
	GOPATH=$(PWD)/.GOPATH
endif

GO_SOURCES = $(shell find . -name ".GOPATH" -prune -o -type f -name '*.go' -print)
JS_SOURCES = $(shell find app/js -type f -name '*.js' -print)
HTML_SOURCES = $(shell find app/html -type f -name '*.html' -print)
CSS_SOURCES = $(shell find app/less -type f -name "*.less" -print)
WEB_LIBS = $(shell find app/lib -type f -type f -print)

DISTS  = dist/js/common.js dist/js/edit.js dist/js/attach.js
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
	go test ./backend/...
	./$(BIN_NAME) -D serve
	#./$(BIN_NAME) -D -w ~/src/shipdock/shipdock serve -o front-page=readme.md
auto-run:
	while true; do \
		make .sources | entr -rd make run ;  \
		echo "hit ^C again to quit" && sleep 1  \
	; done
reset:
	ps -e | grep make | grep -v grep | awk '{print $$1}' | xargs kill

## Binary build
$(BIN_NAME).bin: $(GO_SOURCES) $(GOPATH)/src/$(IMPORT_PATH)
	go get -v -d $(IMPORT_PATH)            # can replace with glide
	go build \
		-ldflags "-X main.Version=$(VERSION)" \
		-o $(BIN_NAME).bin .
	@echo Build DONE

$(BIN_NAME): $(BIN_NAME).bin $(DISTS)
	cp $(BIN_NAME).bin $(BIN_NAME).tmp
	rice append -v --exec $(BIN_NAME).tmp \
		-i $(IMPORT_PATH)/server  \
		-i $(IMPORT_PATH)/server/renderer \
		-i $(IMPORT_PATH)/backend/config
	mv $(BIN_NAME).tmp $(BIN_NAME)
	@echo Embed resources DONE

## Web dist
dist/js/%.js: $(JS_SOURCES)
	traceur \
		--async-functions \
		--modules inline \
		--source-maps=file \
		--inline $(@:dist/js/%.js=app/js/%.js) \
		--out $@
dist/html/%.html: app/html/%.html
	@mkdir -p $(basename $@)
	cp $< $@
dist/css/common.css: $(CSS_SOURCES)
	lessc app/less/main.less $@
dist/lib/%: app/lib/%
	@mkdir -p $(basename $@)
	cp $< $@

tools:
	npm install -g traceur
	npm install -g less
	go get github.com/GeertJohan/go.rice/rice
clean:
	rm -rf dist/ vendor/ $(BIN_NAME) $(BIN_NAME).bin $(BIN_NAME).tmp
	go clean

$(GOPATH)/src/$(IMPORT_PATH):
	@echo "make symbolic link on $(GOPATH)/src/$(IMPORT_PATH)..."
	@mkdir -p $(dir $(GOPATH)/src/$(IMPORT_PATH))
	ln -s $(PWD) $(GOPATH)/src/$(IMPORT_PATH)

.PHONY: .sources run auto-run reset tools clean
