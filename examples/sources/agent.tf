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

resource "mezmo_agent_source" "source1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My agent source"
  description = "This receives data from my webhook"
}

resource "mezmo_agent_source" "shared_source" {
  pipeline_id      = mezmo_pipeline.pipeline1.id
  title            = "A shared agent source"
  description      = "This source uses the same data as source1"
  gateway_route_id = mezmo_agent_source.source1.gateway_route_id
}
