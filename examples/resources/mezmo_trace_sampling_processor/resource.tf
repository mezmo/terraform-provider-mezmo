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

resource "mezmo_open_telemetry_traces_source" "source1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My Open Telemetry Traces source"
  description = "This receives data from Open Telemetry"
}

resource "mezmo_trace_sampling_processor" "processor1" {
  pipeline_id          = mezmo_pipeline.pipeline1.id
  title                = "Tail Based Sampler"
  description          = "Sample the traces based on the tail of the trace"
  inputs               = [mezmo_open_telemetry_traces_source.source1.id]
  sample_type          = "tail"
  trace_id_field       = ".trace_id"
  parent_span_id_field = ".parent_span_id"
  conditionals = [
    {
      rate = 1,
      conditional = {
        expressions = [
          {
            field        = ".status"
            operator     = "greater_or_equal"
            value_number = "500"
          }
        ]
      }
    },
    {
      rate = 10,
      conditional = {
        expressions = [
          {
            field        = ".status"
            operator     = "greater_or_equal"
            value_number = "400"
          },
          {
            field        = ".status"
            operator     = "less"
            value_number = "500"
          }
        ]
        logical_operation = "AND"
      }
    },
    {
      rate = 100,
      conditional = {
        expressions = [
          {
            field        = ".status"
            operator     = "greater_or_equal"
            value_number = "200"
          },
          {
            field        = ".status"
            operator     = "less"
            value_number = "300"
          }
        ]
        logical_operation = "AND"
      }
    },
  ]
}

resource "mezmo_trace_sampling_processor" "processor2" {
  pipeline_id    = mezmo_pipeline.pipeline1.id
  title          = "Head Based Sampler"
  description    = "Keep 1% of traces"
  inputs         = [mezmo_open_telemetry_traces_source.source1.id]
  sample_type    = "head"
  trace_id_field = ".trace_id"
  rate           = 100
}

resource "mezmo_blackhole_destination" "destination1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "Tail based logs destination"
  ack_enabled = false
  inputs      = [mezmo_trace_sampling_processor.processor1]
}

resource "mezmo_blackhole_destination" "destination2" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "Head based logs destination"
  ack_enabled = false
  inputs      = [mezmo_trace_sampling_processor.processor2]
}
