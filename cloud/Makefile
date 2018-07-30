.DEFAULT_GOAL := deploy

export COMPONENT_NAME     ?= bubbles

export TF_VAR_domain      ?= dev.kuubernetes.delivery
export AWS_PROFILE        ?= default
export AWS_DEFAULT_REGION ?= us-east-2
export TF_LOG             ?= info
export TF_LOG_PATH        ?= .terraform/$(TF_VAR_domain).log
export TF_OPTS            ?= -no-color
export TF_UPDATE          ?= -update

STATE_BUCKET ?= terraform.agilestacks.com
STATE_REGION ?= us-east-1

terraform ?= terraform-v0.11

deploy: init plan apply

init:
	@rm -rf .terraform/$(TF_VAR_domain)
	@mkdir -p .terraform/$(TF_VAR_domain)/.terraform
	@cp *.tf .terraform/$(TF_VAR_domain)/
	cd .terraform/$(TF_VAR_domain) && $(terraform) init -get=true -input=false $(TF_OPTS) \
		-backend=true \
		-backend-config="bucket=$(STATE_BUCKET)" \
		-backend-config="region=$(STATE_REGION)" \
		-backend-config="key=$(TF_VAR_domain)/bubbles/$(COMPONENT_NAME)/terraform.tfstate" \
		-backend-config="profile=$(AWS_PROFILE)" \
		.
.PHONY: init

get:
	cd .terraform/$(TF_VAR_domain) && \
		$(terraform) get $(TF_UPDATE)
.PHONY: get

plan:
	cd .terraform/$(TF_VAR_domain) && \
		$(terraform) plan $(TF_OPTS) -refresh=true -module-depth=-1 -out=terraform.tfplan
.PHONY: plan

apply:
	cd .terraform/$(TF_VAR_domain) && \
		$(terraform) apply $(TF_OPTS) -Xshadow=false terraform.tfplan
.PHONY: apply

undeploy: init
	cd .terraform/$(TF_VAR_domain) && \
		$(terraform) plan $(TF_OPTS) -destroy -refresh=true -module-depth=-1 -out=terraform.tfplan && \
		$(terraform) apply $(TF_OPTS) -Xshadow=false terraform.tfplan
.PHONY: undeploy
