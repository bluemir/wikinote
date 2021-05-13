
dev-run:
	@$(MAKE) build/tools/entr
	while true; do \
		$(MAKE) .watched_sources | \
		entr -rd $(MAKE) test run ;  \
		echo "hit ^C again to quit" && sleep 1 \
	; done

reset:
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
