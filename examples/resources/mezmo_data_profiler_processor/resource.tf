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

resource "mezmo_data_profiler_processor" "processor1" {
  pipeline_id  = mezmo_pipeline.pipeline1.id
  title        = "Data profiler"
  description  = "Profile all the data!"
  inputs       = [mezmo_http_source.curl.id]
  app_fields   = [".app", ".container"]
  host_fields  = [".host", ".hostname"]
  level_fields = [".level", ".log_level"]
  line_fields  = [".line", ".message"]
  label_fields = [".labels"]
}
