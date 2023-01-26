# versions.tf

terraform {
  required_version = ">= 1.3.7"

  backend "remote" {
    hostname     = "app.terraform.io"
    organization = "witsoft"
    workspaces { prefix = "index-enron-mail-zincsearch-" }
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.51.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = ">= 2.3.0"
    }
  }
}
