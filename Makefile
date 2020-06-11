.DEFAULT_GOAL := build

OS := $(shell uname -s | tr A-Z a-z)

export GOBIN := $(abspath .)/bin/$(OS)

export AWS_DEFAULT_REGION ?= us-east-2

DOMAIN_NAME    ?= test.dev.superhub.io
COMPONENT_NAME ?= bubbles
LOCAL_IMAGE    ?= agilestacks/bubbles
REGISTRY       ?= $(subst https://,,$(lastword $(shell aws ecr get-login --region $(AWS_DEFAULT_REGION))))
IMAGE          ?= $(REGISTRY)/agilestacks/$(DOMAIN_NAME)/bubbles
IMAGE_VERSION  ?= $(shell git rev-parse HEAD | colrm 7)
NAMESPACE      ?= automation-hub

kubectl ?= kubectl --context="$(DOMAIN_NAME)" --namespace="$(NAMESPACE)"
docker  ?= docker
aws     ?= aws

get:
	go get github.com/agilestacks/bubbles/cmd/bubbles
.PHONY: get

run: get
	$(GOBIN)/bubbles -trace
.PHONY: run

fmt:
	go fmt github.com/agilestacks/bubbles/...
.PHONY: fmt

vet:
	go vet -composites=false github.com/agilestacks/bubbles/...
.PHONY: vet

deploy: build ecr-login push kubernetes

clean:
	@rm -f bubbles bin/bubbles
	@rm -rf bin/darwin bin/linux
.PHONY: clean

build:
	$(docker) build -t $(LOCAL_IMAGE):$(IMAGE_VERSION) -t $(LOCAL_IMAGE):latest .
.PHONY: build

ecr-login:
	$(aws) ecr get-login --no-include-email --region $(AWS_DEFAULT_REGION) | $(SHELL) -
.PHONY: ecr-login

push:
	$(docker) tag $(LOCAL_IMAGE):$(IMAGE_VERSION) $(IMAGE):$(IMAGE_VERSION)
	$(docker) tag $(LOCAL_IMAGE):$(IMAGE_VERSION) $(IMAGE):latest
	$(docker) push $(IMAGE):$(IMAGE_VERSION)
	$(docker) push $(IMAGE):latest
.PHONY: push

kubernetes:
	$(kubectl) apply -f templates/namespace.yaml
	$(kubectl) apply -f templates/service.yaml
	$(kubectl) apply -f templates/secret.yaml
	$(kubectl) apply -f templates/deployment.yaml
.PHONY: kubernetes

undeploy:
	-$(kubectl) delete -f templates/deployment.yaml
	-$(kubectl) delete -f templates/secret.yaml
	-$(kubectl) delete -f templates/service.yaml
.PHONY: undeploy
