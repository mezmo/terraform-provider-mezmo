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

resource "mezmo_kafka_sink" "simple_sink" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My kafka sink"
  description = "This represents a kafka sink"
  inputs      = [mezmo_demo_source.source1.id]

  brokers = [{
    host = "brokers.kafka.com"
    port = 9092
  }]

  topic       = "my-topic"
  compression = "none"
  encoding    = "json"
}

resource "mezmo_kafka_sink" "sink_with_sasl" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My kafka sink"
  description = "This represents a kafka sink using SASL authentication"
  inputs      = [mezmo_demo_source.source1.id]

  brokers = [{
    host = "brokers.kafka.com"
    port = 9092
  }]

  topic       = "my-topic"
  compression = "none"
  encoding    = "json"

  sasl = {
    username  = "my-username"
    password  = "my-password"
    mechanism = "PLAIN"
  }
}

