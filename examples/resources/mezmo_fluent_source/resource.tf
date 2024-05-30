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

resource "mezmo_http_source" "webhook" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My HTTP source"
  description = "This receives Fluent data from my webhook"
  decoding    = "json"
}

resource "mezmo_fluent_source" "direct" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "A Fluent source"
  description = "This receives data directly from Fluent Bit"
}

resource "mezmo_fluent_source" "webhook" {
  pipeline_id      = mezmo_pipeline.pipeline1.id
  title            = "A shared source"
  description      = "This Fluent source uses the webhook data"
  shared_source_id = mezmo_http_source.webhook.shared_source_id
}
