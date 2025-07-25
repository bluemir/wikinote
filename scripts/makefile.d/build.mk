##@ Build
## Go Sources
GO_SOURCES = $(shell find . -name "vendor"  -prune -o \
                            -type f -name "*.go" -print)

build/docker-image: $(GO_SOURCES)

.PHONY: build
build: build/$(APP_NAME) ## Build web app

.PHONY: test
test: fmt vet ## Run test
	go test -v ./...

build/$(APP_NAME): $(GO_SOURCES) $(MAKEFILE_LIST) fmt vet
	@$(MAKE) build/tools/go
	mkdir -p build
	go build -v \
		-trimpath \
		-ldflags "\
			-X '$(IMPORT_PATH)/internal/buildinfo.AppName=$(APP_NAME)' \
			-X '$(IMPORT_PATH)/internal/buildinfo.Version=$(VERSION)' \
			-X '$(IMPORT_PATH)/internal/buildinfo.BuildTime=$(shell go run scripts/tools/date/main.go)' \
		" \
		$(OPTIONAL_BUILD_ARGS) \
		-o $@ .
		#-tags embed \

build-tools: build/tools/go
build/tools/go:
	@which $(notdir $@) || echo "see https://golang.org/doc/install"

.PHONY: fmt
fmt: ## Run go fmt against code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code
	@mkdir -p build/static
	go vet ./...
