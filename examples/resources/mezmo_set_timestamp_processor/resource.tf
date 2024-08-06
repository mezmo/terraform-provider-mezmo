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

resource "mezmo_set_timestamp_processor" "one_field" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "Timestamp Pattern"
  description = "Setting a strftime format.  Will be interpreted as Custom and the pattern will be set as the custom pattern."
  inputs      = [mezmo_http_source.curl.id]
  parsers = [
    {
      field            = ".field1"
      timestamp_format = "%Y-%m-%dT%H:%M:%S"
    },
  ]
}
resource "mezmo_set_timestamp_processor" "two_fields" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "Two fields and patterns"
  description = "Adding a second strftime compatible format.  First set to match will set the event's time."
  inputs      = [mezmo_http_source.curl.id]
  parsers = [
    {
      field            = ".field2"
      timestamp_format = "%Y-%m-%dT%H:%M:%S"
    },
    {
      field            = ".field3"
      timestamp_format = "%m/%d/%Y::%H:%M"
    },
  ]
}