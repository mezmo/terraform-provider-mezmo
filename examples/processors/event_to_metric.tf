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

resource "mezmo_event_to_metric_processor" "processor1" {
  pipeline_id     = mezmo_pipeline.pipeline1.id
  title           = "My event to metric processor"
  description     = "Transform my events into metrics"
  inputs          = [mezmo_http_source.curl.id]
  metric_name     = "my_metric"
  metric_type     = "gauge"
  metric_kind     = "incremental"
  namespace_field = ".namespace"
  value_field     = ".connections"
  tags = [
    {
      name       = "my_tag"
      value_type = "value"
      value      = "tag_value"
    }
  ]
}
