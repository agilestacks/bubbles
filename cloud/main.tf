terraform {
  required_version = ">= 0.11.3"
  backend "s3" {}
}

provider "aws" {
  version = "1.35.0"
}

module "ecr" {
  source = "github.com/agilestacks/terraform-modules//ecr"
  name   = "agilestacks/${var.domain}/bubbles"
}
