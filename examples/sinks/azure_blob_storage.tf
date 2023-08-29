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

resource "mezmo_azure_blob_storage_sink" "sink1" {
  pipeline_id       = mezmo_pipeline.pipeline1.id
  title             = "My sink"
  description       = "Send logs to Azure Blob Storage"
  inputs            = [mezmo_demo_source.source1.id]
  connection_string = "AccountName=mylogstorage;AccountKey=storageaccountkeybase64encoded;EndpointSuffix=core.windows.net"
  container_name    = "my-logs"
  compression       = "gzip"
}
