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

resource "mezmo_demo_source" "source1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My source"
  format      = "apache_common"
}

resource "mezmo_route_processor" "processor1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My processor"
  description = "This processor routes error logs"
  inputs      = [mezmo_demo_source.source1.id]
  conditionals = [
    {
      label = "error logs"
      expressions = [
        {
          field        = ".status"
          operator     = "greater_or_equal"
          value_number = 300
        },
        {
          field        = ".level"
          operator     = "greater_or_equal"
          value_string = "info"
        }
      ]
      logical_operation = "OR"
    }
  ]
}

resource "mezmo_route_processor" "processor2" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My processor"
  description = "This processor routes logs"
  inputs      = [mezmo_demo_source.source1.id]
  conditionals = [
    {
      expressions_group = [
        {
          expressions = [
            {
              field        = ".label"
              operator     = "equal"
              value_string = "account"
            },
            {
              field        = ".app_name"
              operator     = "ends_with"
              value_string = "service"
            },
          ]
        },
        {
          expressions_group = [
            {
              expressions = [
                {
                  field        = ".level"
                  operator     = "greater_or_equal"
                  value_number = 300
                },
                {
                  field        = ".tag"
                  operator     = "contains"
                  value_string = "error"
                }
              ]
              logical_operation = "OR"
            }
          ]
        }
      ]
      label = "error logs"
    }
  ]
}

resource "mezmo_route_processor" "processor3" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My processor"
  description = "This processor routes logs"
  inputs      = [mezmo_demo_source.source1.id]
  conditionals = [
    {
      expressions_group = [
        {
          expressions = [
            {
              field        = ".label"
              operator     = "equal"
              value_string = "account"
            },
            {
              field        = ".app_name"
              operator     = "ends_with"
              value_string = "service"
            },
          ]
        },
        {
          expressions_group = [
            {
              expressions = [
                {
                  field        = ".level"
                  operator     = "greater_or_equal"
                  value_number = 300
                },
                {
                  field        = ".tag"
                  operator     = "contains"
                  value_string = "error"
                }
              ]
              logical_operation = "OR"
            }
          ]
        }
      ]
      label = "error logs"
    },
    {
      expressions = [
        {
          field        = ".status"
          operator     = "equal"
          value_number = 503
        }
      ]
      label = "503 logs"
    }
  ]
}

resource "mezmo_logs_destination" "destination1" {
  pipeline_id   = mezmo_pipeline.pipeline1.id
  title         = "My destination"
  description   = "Send logs to Mezmo Log Analysis"
  inputs        = [mezmo_route_processor.processor3.conditionals.0.output_name]
  ingestion_key = var.my_ingestion_key
}

resource "mezmo_blackhole_destination" "destination2" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My destination"
  description = "Trash the data without acking"
  ack_enabled = false
  inputs      = [mezmo_route_processor.processor3.conditionals.1.output_name]
}

resource "mezmo_blackhole_destination" "destination3" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My destination"
  description = "Send unmatched data to blackhole"
  ack_enabled = false
  inputs      = [mezmo_route_processor.processor3.unmatched]
}
