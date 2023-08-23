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

resource "mezmo_http_source" "curl" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My data stream"
  description = "Send Curl data to the pipeline point of entry URL"
  decoding    = "json"
}

resource "mezmo_sample_transform" "transform1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My sample transform"
  description = "Let's sample all data"
  inputs      = [mezmo_http_source.curl.id]
  rate        = 10
}

resource "mezmo_sample_transform" "transform2" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My second sample transform"
  description = "Let's sample some data while keeping other events intact"
  inputs      = [mezmo_http_source.curl.id]
  rate        = 100
  always_include = {
    field        = ".my_app_id"
    operator     = "greater"
    value_number = 10
  }
}
