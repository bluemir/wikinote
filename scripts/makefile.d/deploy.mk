##@ Deployments
.PHONY: deploy
deploy: build/docker-image.pushed ## Deploy webapp
	#@$(MAKE) runtime/tools/kubectl
	# deploy code
	# example:
	#   cat deploy.yaml | DEPLOY_IMAGE=$(shell cat $<) envsubst | kubectl apply -f -
	#   kubectl kustermize deploy | kubectl apply -f -

install: build/$(APP_NAME)-$(VERSION) ## Install binary to local system
	cp $< $(HOME)/.local/bin/$(APP_NAME)-$(VERSION)
	ln -sf $(APP_NAME)-$(VERSION) $(HOME)/.local/bin/$(APP_NAME)

tools: runtime/tools/kubectl
runtime/tools/kubectl:
	@which $(notdir $@) || (./scripts/tools/install/kubectl.sh)
	#install kubectl. https://kubernetes.io/docs/tasks/tools/