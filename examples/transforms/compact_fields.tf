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

resource "mezmo_compact_fields_transform" "transform1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My compact fields transform"
  description = "Get those null values outta here!"
  inputs      = [mezmo_http_source.curl.id]
  fields      = [".root_level_field", ".nested.field"]
}
