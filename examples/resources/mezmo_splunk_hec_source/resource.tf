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

resource "mezmo_splunk_hec_source" "splunk_source" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "A Splunk HEC Source"
  description = "This source receives data from Splunk HTTP Event Collector"
}
