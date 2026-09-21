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

resource "mezmo_open_telemetry_combined_source" "source1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My open-telemetry source"
  description = "Logs, metrics & traces sent by OpenTelemetry senders using my account's ingestion key."
}

resource "mezmo_logs_destination" "logs_destination" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "OTel logs destination"
  inputs      = [mezmo_open_telemetry_combined_source.source1.logs_output_id]
}

resource "mezmo_prometheus_remote_write_destination" "metrics_destination" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "OTel metrics destination"
  inputs      = [mezmo_open_telemetry_combined_source.source1.metrics_output_id]
  endpoint    = "https://example.org:8080/api/v1/push"
}

resource "mezmo_blackhole_destination" "traces_destination" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "OTel traces destination"
  inputs      = [mezmo_open_telemetry_combined_source.source1.traces_output_id]
}
