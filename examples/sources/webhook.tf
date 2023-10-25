terraform {
  required_providers {
    mezmo = {
      source = "registry.terraform.io/mezmo/mezmo"
    }
  }
  required_version = ">= 1.1.0"
}

provider "mezmo" {
  auth_key = "my secret"
}

resource "mezmo_pipeline" "pipeline1" {
  title = "My pipeline"
}

resource "mezmo_webhook_source" "my_webhook" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My webhook source"
  description = "This is a source made from a webhook call"
  signing_key = "sshhhh"
}
