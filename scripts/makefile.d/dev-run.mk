##@ Run


.PHONY: run
run: build/$(APP_NAME) ## Run web app
	$< -vv server --wiki-path=runtime --volatile-database

.PHONY: dev-run
dev-run: ## Run dev server. If detect file change, automatically rebuild&restart server
	@$(MAKE) build/tools/watcher
	watcher \
		--include "go.mod" \
		--include "go.sum" \
		--include "**.go" \
		--include "internal/**" \
		--include "package.json" \
		--include "yarn.lock" \
		--include "assets/**" \
		--include "api/proto/**" \
		--include "Makefile" \
		--include "scripts/makefile.d/*.mk" \
		--include "runtime/init-data.yaml" \
		--include "runtime/config.yaml" \
		--exclude "build/**" \
		--exclude "**.sw*" \
		--exclude "internal/swagger/**" \
		--exclude "assets/js/index.js" \
		--debounce=1s \
		-- \
	$(MAKE) test run


.PHONY: reset
reset: ## Kill all make process. Use when dev-run stuck.
	ps -e | grep make | grep -v grep | awk '{print $$1}' | xargs kill

build/tools/entr:
	@which $(notdir $@) || (echo "see http://eradman.com/entrproject")

tools: build/tools/watcher
build/tools/watcher:
	@which $(notdir $@) || (./scripts/tools/install/go-tool.sh github.com/bluemir/watcher)

