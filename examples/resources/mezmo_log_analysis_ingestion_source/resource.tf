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

resource "mezmo_log_analysis_ingestion_source" "source1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My log-analysis-ingestion source"
  description = "Logs sent to Mezmo's Log Analysis endpoint, redirected to my pipeline first."
}
