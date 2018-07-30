terraform {
  required_version = ">= 0.11.3"
  backend "s3" {}
}

provider "aws" {}

module "ecr" {
  source = "github.com/agilestacks/terraform-modules//ecr"
  name   = "agilestacks/${var.domain}/bubbles"
}
