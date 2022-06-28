terraform {
  backend "remote" {
    organization = "numary"

    workspaces {
      prefix = "app-webhooks-"
    }
  }
}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "1.3.1"
    }
  }
}

provider "aws" {
  region = local.region
  default_tags {
    tags = {
      Environment = var.env
      App         = "webhooks"
    }
  }
}

locals {
  region = "eu-west-1"
}

variable "env" {}
variable "vpc_id" {}
variable "vpc_cidr" {}
variable "subnet_a" {}
variable "subnet_b" {}
variable "subnet_c" {}
variable "app_env_name" {}
variable "atlas_project_id" {}
