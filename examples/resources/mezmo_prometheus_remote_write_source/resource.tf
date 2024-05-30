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

resource "mezmo_prometheus_remote_write_source" "source1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My Prometheus Remote Write source"
  description = "This receives data from prometheus"
}

resource "mezmo_prometheus_remote_write_source" "shared_source" {
  pipeline_id      = mezmo_pipeline.pipeline1.id
  title            = "A shared Prometheus Remote Write source"
  description      = "This source uses the same data as source1"
  shared_source_id = mezmo_prometheus_remote_write_source.source1.shared_source_id
}
