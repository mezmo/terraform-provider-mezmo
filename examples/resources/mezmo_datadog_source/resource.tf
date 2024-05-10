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

resource "mezmo_datadog_source" "datadog1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "A Datadog source for Mezmo pipelines"
  description = "This can be a parent source"
}

resource "mezmo_datadog_source" "shared_datadog_source" {
  pipeline_id      = mezmo_pipeline.pipeline1.id
  title            = "A shared Datadog source"
  description      = "This source uses the same data as datadog1"
  gateway_route_id = mezmo_datadog_source.datadog1.gateway_route_id
}
