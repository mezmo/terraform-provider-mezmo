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

resource "mezmo_reduce_processor" "processor1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My processor"
  description = "This processor reduces data"
  inputs      = [mezmo_demo_source.source1.id]
  group_by    = [".field1.id"]
  date_formats = [
    {
      field  = ".datetime"
      format = "%d/%m/%Y:%T"
    }
  ]
  merge_strategies = [
    {
      field    = ".method"
      strategy = "array"
    }
  ]
  flush_condition = {
    when = "starts_when"
    conditional = {
      expressions = [
        {
          field        = ".status"
          operator     = "equal"
          value_number = 200
        },
        {
          field        = ".status"
          operator     = "equal"
          value_number = 201
        }
      ],
      logical_operation = "OR"
    }
  }
}
