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

resource "mezmo_demo_source" "source1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My source"
  format      = "apache_common"
}

resource "mezmo_filter_processor" "processor1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My processor"
  description = "This processor filters logs"
  inputs      = [mezmo_demo_source.source1.id]
  conditionals = {
    expressions = [
      {
        field        = ".status"
        operator     = "greater_or_equal"
        value_number = 300
      },
      {
        field        = ".level"
        operator     = "contains"
        value_string = "info"
      }
    ]
    logical_operation = "OR"
  }
}
