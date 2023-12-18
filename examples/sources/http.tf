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

resource "mezmo_http_source" "source1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My HTTP source"
  description = "This receives data from HTTP clients"
  decoding    = "json"
}

resource "mezmo_http_source" "shared_source" {
  pipeline_id      = mezmo_pipeline.pipeline1.id
  title            = "A shared HTTP source"
  description      = "This source uses the same data as source1"
  decoding         = "json"
  gateway_route_id = mezmo_http_source.source1.gateway_route_id
}
