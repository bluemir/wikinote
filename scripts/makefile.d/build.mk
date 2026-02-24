##@ Build
## Go Sources
GO_SOURCES = $(shell find . -name "vendor"  -prune -o \
                            -type f -name "*.go" -print)

build/docker-image: $(GO_SOURCES)

.PHONY: build
build: build/$(APP_NAME) ## Build web app

.PHONY: test
test: fmt vet ## Run test
	go test -trimpath ./...

build/$(APP_NAME): $(GO_SOURCES) $(MAKEFILE_LIST) | fmt vet test runtime/tools/go
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

.PHONY: prod
prod: build/$(APP_NAME)-$(VERSION) ## Build web app (prod mode, minified + embedded)

build/$(APP_NAME)-$(VERSION): $(GO_SOURCES) $(MAKEFILE_LIST) | fmt vet gen test runtime/tools/go
	mkdir -p build
	go build -v \
		-trimpath \
		-tags prod \
		-ldflags "\
			-X '$(IMPORT_PATH)/internal/buildinfo.AppName=$(APP_NAME)' \
			-X '$(IMPORT_PATH)/internal/buildinfo.Version=$(VERSION)' \
			-X '$(IMPORT_PATH)/internal/buildinfo.BuildTime=$(shell go run scripts/tools/date/main.go)' \
		" \
		$(OPTIONAL_BUILD_ARGS) \
		-o $@ .

.PHONY: gen
gen: runtime/tools/esbuild ## Run go generate
	PATH=$(shell pwd)/runtime/tools:$(PATH) go generate -tags prod -x ./...

build-tools: runtime/tools/go runtime/tools/govulncheck runtime/tools/gosec
runtime/tools/go:
	@which $(notdir $@) || echo "see https://golang.org/doc/install"

.PHONY: vulncheck
vulncheck: runtime/tools/govulncheck ## Run govulncheck
	govulncheck ./...

runtime/tools/govulncheck:
	GOBIN=$(shell pwd)/runtime/tools go install golang.org/x/vuln/cmd/govulncheck@latest

.PHONY: sec
sec: runtime/tools/gosec ## Run gosec
	gosec ./...

runtime/tools/gosec:
	GOBIN=$(shell pwd)/runtime/tools go install github.com/securego/gosec/v2/cmd/gosec@latest

.PHONY: fmt
fmt: ## Run go fmt against code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code
	@mkdir -p build/static
	go vet ./...
