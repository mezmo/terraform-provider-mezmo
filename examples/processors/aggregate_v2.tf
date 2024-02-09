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

resource "mezmo_aggregate_v2_processor" "processor1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My aggregate v2 processor"
  description = "Aggregate my metrics via tumbling window"
  method      = "tumbling"
  interval    = 3600
}

resource "mezmo_aggregate_v2_processor" "processor2" {
  pipeline_id     = mezmo_pipeline.pipeline1.id
  title           = "My aggregate v2 processor"
  description     = "Aggregate my metrics via sliding window"
  method          = "sliding"
  window_duration = 3600
  strategy        = "AVG"
  conditional = {
    expressions = [
      {
        field        = ".tags.host"
        operator     = "contains"
        value_string = "internal"
      },
      {
        field        = ".tags.app"
        operator     = "contains"
        value_string = "my_app"
      }
    ]
    logical_operation = "OR"
  }
}
