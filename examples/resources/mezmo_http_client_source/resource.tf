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

resource "mezmo_http_client_source" "source1" {
  pipeline_id          = mezmo_pipeline.pipeline1.id
  title                = "My HTTP client source"
  description          = "This polls an HTTP endpoint for data"
  endpoint             = "https://example.com/data"
  method               = "GET"
  scrape_interval_secs = 30
  scrape_timeout_secs  = 10
  decoding             = "json"

  auth = {
    strategy = "bearer"
    token    = "my-secret-token"
  }

  headers = {
    "X-Custom-Header" = "custom-value"
  }
}
