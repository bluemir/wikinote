##@ Docker
DOCKER_IMAGE_NAME=$(shell echo $(APP_NAME)| tr A-Z a-z)
DOCKER?=$(shell scripts/tools/find-docker-alt.sh)

.PHONY: docker
docker: build/docker-image ## Build docker image

build/docker-image: Dockerfile $(MAKEFILE_LIST) | runtime/tools/docker
	@mkdir -p $(dir $@)
	$(DOCKER) build \
		--build-arg VERSION=$(VERSION) \
		-t $(DOCKER_IMAGE_NAME):$(VERSION) \
		-f $< .
	echo $(DOCKER_IMAGE_NAME):$(VERSION) > $@

.PHONY: docker-push
docker-push: build/docker-image.pushed ## Push docker image

build/docker-image.pushed: build/docker-image | runtime/tools/docker
	@mkdir -p $(dir $@)
	$(DOCKER) push $(shell cat $<)
	echo $(shell cat $<) > $@

.PHONY:  docker-run
docker-run: build/docker-image ## Run docker container
	$(DOCKER) run -it --rm \
	-v $(PWD)/runtime:/var/run/data:Z \
	$(shell cat $<) \
	$(APP_NAME) -vv server -w /var/run/data

tools: runtime/tools/docker
runtime/tools/docker:
	@which $(DOCKER)
