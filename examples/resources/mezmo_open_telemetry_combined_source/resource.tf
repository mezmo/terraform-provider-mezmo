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
