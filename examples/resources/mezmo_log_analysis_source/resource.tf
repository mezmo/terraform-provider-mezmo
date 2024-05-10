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

resource "mezmo_log_analysis_source" "source1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My log-analysis source"
  description = "This automatically attaches to my Log Analysis account and streams data"
}
