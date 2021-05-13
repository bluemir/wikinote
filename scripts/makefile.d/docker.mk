DOCKER_IMAGE_NAME=$(shell echo $(APP_NAME)| tr A-Z a-z)

docker: build/docker-image

build/docker-image: Dockerfile $(MAKEFILE_LIST)
	@$(MAKE) build/tools/docker
	@mkdir -p $(dir $@)
	docker build \
		--build-arg VERSION=$(VERSION) \
		-t $(DOCKER_IMAGE_NAME):$(VERSION) \
		-f $< .
	echo $(DOCKER_IMAGE_NAME):$(VERSION) > $@

docker-push: build/docker-image.pushed

build/docker-image.pushed: build/docker-image
	@$(MAKE) build/tools/docker
	@mkdir -p $(dir $@)
	docker push $(shell cat $<)
	echo $(shell cat $<) > $@

docker-run: build/docker-image
	docker run -it --rm -v $(PWD)/runtime:/runtime -w=/runtime $(shell cat $<) $(APP_NAME) -vvvv server

.watched_sources: Dockerfile

tools: build/tools/docker
build/tools/docker:
	@which $(notdir $@)

.PHONY: docker docker-push docker-run
