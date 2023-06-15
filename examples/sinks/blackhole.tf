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
  pipeline = mezmo_pipeline.pipeline1.id
  title    = "My source"
  format   = "nginx"
}

resource "mezmo_blackhole_sink" "sink1" {
  pipeline = mezmo_pipeline.pipeline1.id
  title    = "My sink"
  inputs   = [mezmo_demo_source.source1.id]
}
