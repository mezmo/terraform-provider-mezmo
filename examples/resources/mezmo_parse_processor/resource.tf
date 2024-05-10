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

resource "mezmo_parse_processor" "error_logs" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "Apache parser"
  description = "Parse apache logs"
  inputs      = [mezmo_http_source.curl.id]
  field       = ".log"
  parser      = "apache_log"
  apache_log_options = {
    format           = "error"
    timestamp_format = "%Y/%m/%d %H:%M:%S"
  }
}

resource "mezmo_parse_processor" "retrieve_timestamp" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "Apache parser"
  description = "Parse apache logs"
  inputs      = [mezmo_http_source.curl.id]
  field       = ".log"
  parser      = "timestamp_parser"
  timestamp_parser_options = {
    format        = "Custom"
    custom_format = "%Y/%m/%d %H:%M:%S"
  }
}

resource "mezmo_parse_processor" "regex" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "Regex parser"
  description = "Parse regex log"
  inputs      = [mezmo_http_source.curl.id]
  field       = ".log"
  parser      = "regex_parser"
  regex_parser_options = {
    pattern        = "^(?P<number>[0-9]*)(?P<word>\\w*)(?P<singlequote>\\')(?P<slash>\\\\?)"
    case_sensitive = false
  }
}
