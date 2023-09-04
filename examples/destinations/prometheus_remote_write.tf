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
  title       = "My source"
  description = "This is the point of entry for our data"
  format      = "nginx"
}

resource "mezmo_prometheus_remote_write_destination" "destination1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "Prometheus Metrics Destination"
  description = "This is a Prometheus destination that uses remote-write method"
  inputs      = [mezmo_demo_source.source1.id]
  endpoint    = "https://example.org:8080/api/v1/push"
}
