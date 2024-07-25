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
  description = "This is some fake data for testing"
  format      = "generic_metrics"
}

resource "mezmo_gcp_cloud_monitoring_destination" "gcp" {
  title            = "GCP Cloud Monitoring"
  description      = "This stores our metrics events in GCP cloud monitoring"
  inputs           = [mezmo_demo_source.source1.id]
  pipeline_id      = mezmo_pipeline.pipeline1.id
  resource_type    = "global2"
  project_id       = "proj456"
  credentials_json = "{}"
  resource_labels = {
    "somekey1"  = "v1"
    "otherkey1" = "v2"
  }
}
