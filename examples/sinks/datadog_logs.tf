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

resource "mezmo_demo_source" "source1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My source"
  description = "This is the point of entry for our data"
  format      = "json"
}

resource "mezmo_datadog-logs_sink" "sink1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My sink"
  description = "Datadog logs sink"
  site        = "us1"
  compression = "none"
  api_key     = "<secret-api-key>"
  inputs      = [mezmo_demo_source.source1.id]
}
