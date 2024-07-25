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
  format      = "nginx"
}

resource "mezmo_gcp_cloud_operations_destination" "gcp" {
  title            = "GCP Cloud Operations"
  description      = "This stores our log events in GCP cloud operations"
  inputs           = [mezmo_demo_source.source1.id]
  pipeline_id      = mezmo_pipeline.pipeline1.id
  resource_type    = "global2"
  project_id       = "proj456"
  log_id           = "log456"
  credentials_json = "{}"
  resource_labels = {
    "somekey1"  = "v1"
    "otherkey1" = "v2"
  }
}
