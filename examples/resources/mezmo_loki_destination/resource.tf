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
  format      = "json"
}

resource "mezmo_loki_destination" "sink1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "test loki sink"
  description = "loki sink description"
  auth = {
    strategy = "basic"
    user     = "username"
    password = "secret-password"
  }
  endpoint = "http://example.com"
  encoding = "json"
  path     = "example/path"
  labels = {
    test_key_0 = "test_value_0"
    test_key_1 = "test_value_1"
  }
}
