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

resource "mezmo_elasticsearch_sink" "sink1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My sink"
  description = "Send logs to an ElasticSearch cluster"
  inputs      = [mezmo_demo_source.source1.id]
  endpoints   = ["https://my.example.com/"]
  auth = {
    strategy = "basic"
    user     = "usr1"
    password = var.my_password
  }
}
