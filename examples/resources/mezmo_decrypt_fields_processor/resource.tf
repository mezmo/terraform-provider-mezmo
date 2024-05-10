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

resource "mezmo_http_source" "curl" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My data stream"
  description = "Send Curl data to the pipeline point of entry URL"
  decoding    = "json"
}

resource "mezmo_decrypt_fields_processor" "processor1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "Decrypt"
  description = "This decrypts an incoming field that will be encrypted"
  inputs      = [mezmo_http_source.curl.id]
  algorithm   = "AES-128-CFB"
  key         = "1111111111111111"
  iv_field    = ".some_iv_field"
  field       = ".this_is_encrypted"
}
