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
  description = "This is the point of entry for our data"
  format      = "nginx"
}

resource "mezmo_splunk_hec_logs_sink" "sink1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My sink"
  description = "Send logs to a Splunk HEC server"
  inputs      = [mezmo_demo_source.source1.id]
  endpoint    = "https://example3.com"
  token       = var.my_splunk_token
  source = {
    value = "my source"
  }
  index = {
    field = ".my_index"
  }
}
