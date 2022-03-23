##@ Build
## Go Sources
GO_SOURCES = $(shell find . -name "vendor"  -prune -o \
                            -type f -name "*.go" -print)

.watched_sources: $(GO_SOURCES) go.mod go.sum
build/docker-image: $(GO_SOURCES)

.PHONY: build
build: build/$(APP_NAME) ## Build web app

.PHONY: test
test: fmt vet ## Run test
	go test -v ./...

build/$(APP_NAME).unpacked: $(GO_SOURCES) $(MAKEFILE_LIST) fmt vet
	@$(MAKE) build/tools/go
	@mkdir -p build
	go build -v \
		-trimpath \
		-ldflags "\
			-X 'main.AppName=$(APP_NAME)' \
			-X 'main.Version=$(VERSION)'  \
		" \
		$(OPTIONAL_BUILD_ARGS) \
		-o $@ main.go

build/$(APP_NAME): build/$(APP_NAME).unpacked $(MAKEFILE_LIST)
	@$(MAKE) build/tools/rice
	@mkdir -p $(dir $@)
	date --rfc-3339=ns > build/static/.time
	cp $< $@.tmp
	rice append -v \
		-i $(IMPORT_PATH)/internal/static \
		--exec $@.tmp
	mv $@.tmp $@

build-tools: build/tools/go build/tools/rice
build/tools/go:
	@which $(notdir $@) || echo "see https://golang.org/doc/install"
build/tools/rice: build/tools/go
	@which $(notdir $@) || (./scripts/go-install-tool.sh github.com/GeertJohan/go.rice/rice)

.PHONY: fmt
fmt: ## Run go fmt against code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code
	go vet ./...
