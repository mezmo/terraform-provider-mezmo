terraform {
  required_providers {
    mezmo = {
      source = "registry.terraform.io/mezmo/mezmo"
    }
  }
  required_version = ">= 1.1.0"
}

variable "my_aws_secret_access_key" {
  type = string
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

resource "mezmo_sqs_destination" "destination1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My destination"
  description = "Send my events to SQS"
  inputs      = [mezmo_demo_source.source1.id]
  compression = "gzip"
  region      = "us-east2"
  queue_url   = "https://sqs.us-east-2.amazonaws.com/123456789012/my-queue"
  auth = {
    access_key_id     = "my_key"
    secret_access_key = var.my_aws_secret_access_key
  }
}
