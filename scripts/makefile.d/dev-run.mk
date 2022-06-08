##@ Run

run: build/$(APP_NAME) ## Run web app
	$< -vv server --admin-user root=1234 --admin-user bluemir=1234 --wiki-path=runtime
dev-run: ## Run dev server. If detect file change, automatically rebuild&restart server
	@$(MAKE) build/tools/entr
	while true; do \
		$(MAKE) .watched_sources | \
		entr -rd $(MAKE) test run ;  \
		echo "hit ^C again to quit" && sleep 1 \
	; done
reset: ## Kill all make process. Use when dev-run stuck.
	ps -e | grep make | grep -v grep | awk '{print $$1}' | xargs kill

## watched_sources
.watched_sources: \
	$(MAKEFILE_LIST)
	@echo $^ | tr " " "\n"

# To add watched resource, just add as depandancy
# example:
#   .watched_sources: Dockerfile

tools: build/tools/entr
build/tools/entr:
	@which $(notdir $@) || (echo "see http://eradman.com/entrproject")

.PHONY: dev-run reset
