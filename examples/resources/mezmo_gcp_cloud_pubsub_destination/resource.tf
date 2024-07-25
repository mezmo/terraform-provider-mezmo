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

resource "mezmo_gcp_cloud_pubsub_destination" "gcp" {
  title            = "GCP Cloud PubSub"
  description      = "This stores our log events in GCP cloud PubSub"
  inputs           = [mezmo_demo_source.source1.id]
  pipeline_id      = mezmo_pipeline.pipeline1.id
  encoding         = "text"
  project_id       = "proj456"
  topic            = "topic456"
  credentials_json = "{}"
}
