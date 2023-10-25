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

resource "mezmo_sample_processor" "processor1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My second sample processor"
  description = "Let's sample some data while keeping other events intact"
  inputs      = [mezmo_http_source.curl.id]
  rate        = 100
  always_include = {
    field          = ".my_app"
    operator       = "starts_with"
    value_string   = "Server"
    case_sensitive = false
  }
}
