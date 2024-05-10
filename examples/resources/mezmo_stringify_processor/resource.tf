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
  format      = "apache_common"
}

resource "mezmo_stringify_processor" "processor1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My processor"
  description = "This processor removes the data we don't want"
  inputs      = [mezmo_demo_source.source1.id]
}
