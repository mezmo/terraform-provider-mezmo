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
  value       = mezmo_access_key.shared.key
  sensitive   = true
  description = "Output for the one-time seen key"
}

resource "mezmo_access_key" "shared" {
  title     = "access key for the shared source"
  source_id = mezmo_shared_source.http.id
}

resource "mezmo_shared_source" "http" {
  title       = "HTTP shared source"
  description = "This http source can be used across different pipelines"
  type        = "http"
}
