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

resource "mezmo_s3_source" "source1" {
  pipeline      = mezmo_pipeline.pipeline1.id
  title         = "My S3 source"
  description   = "This pulls data from S3 using their SQS service"
  region        = "us-east-2"
  sqs_queue_url = "https://hello.com/sqs"
  auth = {
    access_key_id     = "123"
    secret_access_key = "secret123"
  }
}

