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

resource "mezmo_drop_fields_processor" "processor1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My drop fields processor"
  description = "Drop these properties!"
  inputs      = [mezmo_http_source.curl.id]
  fields      = [".drop_this_prop", ".and_also_this"]
}
