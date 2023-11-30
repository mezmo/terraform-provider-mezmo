terraform {
  required_providers {
    mezmo = {
      source = "registry.terraform.io/mezmo/mezmo"
    }
  }
  required_version = ">= 1.1.0"
}

variable "my_ingestion_key" {
  type = string
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

resource "mezmo_parse_sequentially_processor" "processor1" {
  pipeline_id  = mezmo_pipeline.pipeline1.id
  title        = "Apache parser"
  description  = "Parse logs"
  inputs       = [mezmo_http_source.curl.id]
  field        = ".log"
  target_field = ".log_parsed"
  parsers = [
    {
      parser = "apache_log"
      apache_log_options = {
        format           = "error"
        timestamp_format = "%Y/%m/%d %H:%M:%S"
      }
    },
    {
      parser = "json_parser"
    },
    {
      parser = "csv_row"
      csv_parser_options = {
        field_names = ["field_1", "field_2"]
      }
    },
    {
      parser = "regex_parser"
      regex_parser_options = {
        pattern        = "\\d{3}-\\d{2}-\\d{3}"
        case_sensitive = false
        multiline      = true
      }
    }
  ]
}

resource "mezmo_logs_destination" "destination1" {
  pipeline_id   = mezmo_pipeline.pipeline1.id
  title         = "My destination"
  description   = "Send logs to Mezmo Log Analysis"
  inputs        = [mezmo_parse_sequentially_processor.processor1.parsers.0.output_name]
  ingestion_key = var.my_ingestion_key
}

resource "mezmo_blackhole_destination" "destination2" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My destination"
  description = "Trash the data without acking"
  ack_enabled = false
  inputs      = [mezmo_parse_sequentially_processor.processor1.parsers.1.output_name]
}

resource "mezmo_http_destination" "destination3" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "Http desintation"
  description = "Send data to an HTTP destination"
  uri         = "https://example.org"
  inputs      = [mezmo_parse_sequentially_processor.processor1.parsers.2.output_name]
}
