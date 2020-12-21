DOCKER_IMAGE_NAME=$(shell echo $(APP_NAME)| tr A-Z a-z)

docker: build/docker-image

build/docker-image: Dockerfile $(GO_SOURCES) $(HTML_SOURCES) $(JS_SOURCES) $(CSS_SOURCES) $(WEB_LIBS)
	@$(MAKE) build/tools/docker
	@mkdir -p $(dir $@)
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg APP_NAME=$(APP_NAME) \
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

.PHONY: docker docker-push
