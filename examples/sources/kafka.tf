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

# Basic configuration
resource "mezmo_kafka_source" "source1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My Kafka source"
  description = "This receives data from kafka"

  brokers = [{
    host = "brokers.kafka.com"
    port = 9092
  }]

  topics   = ["topic1", "topic2"]
  group_id = "my-group-id"
  decoding = "json"
}

# Configuration with SASL enabled
resource "mezmo_kafka_source" "source2" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My Kafka source"
  description = "This receives data from kafka"

  brokers = [{
    host = "brokers.kafka.com"
    port = 9092
  }]

  topics   = ["topic1", "topic2"]
  group_id = "my-group-id"

  sasl = {
    username  = "my-username"
    password  = "my-password"
    mechanism = "PLAIN"
  }

  decoding    = "json"
}


