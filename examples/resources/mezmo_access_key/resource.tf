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

output "ingestion_key" {
  value       = mezmo_access_key.for_http.key
  sensitive   = true
  description = "The key necessary to ingest data into the source"
}

resource "mezmo_pipeline" "my_pipeline" {
  title = "pipeline"
}
resource "mezmo_http_source" "my_source" {
  pipeline_id = mezmo_pipeline.my_pipeline.id
}
resource "mezmo_access_key" "for_http" {
  title     = "http ingestion key"
  source_id = mezmo_http_source.my_source.id
}
