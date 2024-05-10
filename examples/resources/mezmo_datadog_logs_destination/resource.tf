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

resource "mezmo_datadog_logs_destination" "destination1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My destination"
  description = "Datadog logs destination"
  site        = "us1"
  compression = "gzip"
  api_key     = "<secret-api-key>"
  inputs      = [mezmo_demo_source.source1.id]
}
