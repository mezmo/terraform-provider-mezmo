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
  description = "This receives LogStash data from my webhook that can be shared"
  decoding    = "json"
}

resource "mezmo_logstash_source" "text" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "Logstash direct"
  description = "This source receives text data direct from Logstash"
  format      = "text"
}

resource "mezmo_logstash_source" "shared_source" {
  pipeline_id      = mezmo_pipeline.pipeline1.id
  title            = "Logstash from HTTP"
  description      = "This source uses the same data from the HTTP source (in Logstash format)"
  shared_source_id = mezmo_http_source.webhook.shared_source_id
}
