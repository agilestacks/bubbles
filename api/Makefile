.DEFAULT_GOAL := compile

export GOPATH := $(abspath .)
export GOBIN  := $(GOPATH)/bin/$(shell uname -s | tr A-Z a-z)
export PATH   := $(GOBIN):$(PATH)

export AWS_DEFAULT_REGION ?= us-east-2

DOMAIN_NAME     ?= dev.kubernetes.delivery
COMPONENT_NAME  ?= bubbles
REGISTRY        ?= $(subst https://,,$(lastword $(shell aws ecr get-login --region $(AWS_DEFAULT_REGION))))
IMAGE           ?= $(REGISTRY)/agilestacks/$(DOMAIN_NAME)/bubbles
IMAGE_VERSION   ?= $(shell git rev-parse HEAD | colrm 7)
NAMESPACE       ?= automation-hub

kubectl         ?= kubectl --context="$(DOMAIN_NAME)" --namespace="$(NAMESPACE)"
docker          ?= docker
aws             ?= aws

install:
	@go get -u github.com/kardianos/govendor
.PHONY: install

govendor-list:
	@cd src/bubbles && $(GOBIN)/govendor list
.PHONY: govendor-list

govendor: govendor-list
	@cd src/bubbles && $(GOBIN)/govendor sync
.PHONY: govendor

govendor-add: govendor-list
	@cd src/bubbles && $(GOBIN)/govendor add +e
.PHONY: govendor-add

compile: govendor get

get:
	@go get bubbles
.PHONY: get

run: get
	@$(GOBIN)/bubbles -trace
.PHONY: run

fmt:
	@go fmt bubbles bubbles/api bubbles/config bubbles/flags
.PHONY: fmt

deploy: build ecr-login push kubernetes

clean:
	-@rm -rf .cache pkg bin/darwin bin/linux \
		src/github.com src/golang.org src/gopkg.in \
		src/hub/vendor/github.com src/hub/vendor/golang.org src/hub/vendor/gopkg.in
	-@find src -not -path "*src/bubbles*" -not -path "src" -type d -maxdepth 1 | xargs rm -rf
.PHONY: clean

build:
	@$(docker) build -t $(IMAGE):$(IMAGE_VERSION) .
.PHONY: build

ecr-login:
	$(aws) ecr get-login --no-include-email --region $(AWS_DEFAULT_REGION) | $(SHELL) -
.PHONY: ecr-login

push:
	$(docker) tag  $(IMAGE):$(IMAGE_VERSION) $(IMAGE):latest
	$(docker) push $(IMAGE):$(IMAGE_VERSION)
	$(docker) push $(IMAGE):latest
.PHONY: push

kubernetes:
	$(kubectl) apply -f templates/namespace.yaml
	$(kubectl) apply -f templates/service.yaml
	$(kubectl) apply -f templates/deployment.yaml
.PHONY: kubernetes

undeploy:
	-$(kubectl) delete -f templates/deployment.yaml
	-$(kubectl) delete -f templates/service.yaml
.PHONY: undeploy
