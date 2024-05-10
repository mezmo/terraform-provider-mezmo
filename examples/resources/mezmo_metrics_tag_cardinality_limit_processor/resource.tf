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

resource "mezmo_metrics_tag_cardinality_limit_processor" "processor1" {
  pipeline_id  = mezmo_pipeline.pipeline1.id
  title        = "Tag Limiter"
  description  = "Keeps my sysdig from blowing up"
  inputs       = [mezmo_http_source.curl.id]
  tags         = ["server", "host", "api"]
  exclude_tags = ["ip_address"]
  value_limit  = 50
  action       = "drop_tag"
}
