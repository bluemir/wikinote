##@ Run

run: build/$(APP_NAME) ## Run web app
	$< -vv server --admin-user root=1234 --admin-user bluemir=1234 --wiki-path=runtime
dev-run: ## Run dev server. If detect file change, automatically rebuild&restart server
	@$(MAKE) build/tools/watcher
	while true; do \
		watcher -vv \
			--include "go.mod" \
			--include "go.sum" \
			--include "package.json" \
			--include "yarn.lock" \
			--include "Makefile" \
			--include "scripts/makefile.d/**.mk" \
			--include "web/**" \
			--include "**.go" \
			--include "runtime/.app/config.yaml" \
			-- \
		$(MAKE) test run ;  \
		echo "hit ^C again to quit" && sleep 1 \
	; done
reset: ## Kill all make process. Use when dev-run stuck.
	ps -e | grep make | grep -v grep | awk '{print $$1}' | xargs kill

build/tools/entr:
	@which $(notdir $@) || (echo "see http://eradman.com/entrproject")

tools: build/tools/watcher
build/tools/watcher:
	@which $(notdir $@) || (./scripts/makefile.d/go-install-tool.sh github.com/bluemir/watcher)

.PHONY: dev-run reset
