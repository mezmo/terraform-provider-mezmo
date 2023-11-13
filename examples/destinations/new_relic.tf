terraform {
  required_providers {
    mezmo = {
      source = "registry.terraform.io/mezmo/mezmo"
    }
  }
  required_version = ">= 1.1.0"
}

variable "my_new_relic_account_id" {
  type = string
}

variable "my_license_key" {
  type = string
}
provider "mezmo" {
  auth_key = "my secret"
}

resource "mezmo_pipeline" "pipeline1" {
  title = "My pipeline"
}

resource "mezmo_demo_source" "source1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My source"
  description = "This is the point of entry for our data"
  format      = "nginx"
}

resource "mezmo_new_relic_destination" "destination1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My destination"
  description = "Send logs to New Relic"
  inputs      = [mezmo_demo_source.source1.id]
  account_id  = var.my_new_relic_account_id
  license_key = var.my_license_key
}
