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
  action      = "drop"
  inputs      = [mezmo_demo_source.source1.id]
  conditional = {
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
        negate       = true
      }
    ]
    logical_operation = "OR"
  }
}

resource "mezmo_filter_processor" "complex1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My processor"
  description = "This processor filters logs"
  action      = "drop"
  inputs      = [mezmo_demo_source.source1.id]
  conditional = {
    expressions_group = [
      {
        expressions = [
          {
            field        = ".field",
            operator     = "equal",
            value_string = "info"
          },
          {
            field        = ".field2",
            operator     = "starts_with",
            value_string = "pipeline"
          }
        ],
        logical_operation = "AND"
      },
      {
        expressions_group = [
          {
            expressions = [
              {
                field        = ".field3",
                operator     = "ends_with",
                value_string = "error"
              }
            ]
          },
          {
            expressions = [
              {
                field        = ".field4",
                operator     = "equal",
                value_string = "foo"
              },
              {
                field        = ".field5",
                operator     = "less_or_equal",
                value_string = 1000
              }
            ]
            logical_operation = "OR"
          }
        ]
        logical_operation = "OR"
      }
    ],
    "logical_operation" = "OR"
  }
}
