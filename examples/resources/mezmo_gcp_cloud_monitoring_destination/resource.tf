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
  description = "This is some fake data for testing"
  format      = "generic_metrics"
}

resource "mezmo_gcp_cloud_monitoring_destination" "gcp" {
  title            = "GCP Cloud Monitoring"
  description      = "This stores our metrics events in GCP cloud monitoring"
  inputs           = [mezmo_demo_source.source1.id]
  pipeline_id      = mezmo_pipeline.pipeline1.id
  resource_type    = "global2"
  project_id       = "proj456"
  credentials_json = <<-EOT
            {
              "type": "service_account",
              "project_id": "projid",
              "private_key_id": "12345privatekeyID",
              "private_key": "-----BEGIN PRIVATE KEY-----\nprivatekey\n-----END PRIVATE KEY-----\n",
              "client_email": "serviceacct@proj.iam.gserviceaccount.com",
              "client_id": "000000",
              "auth_uri": "https://accounts.google.com/o/oauth2/auth",
              "token_uri": "https://oauth2.googleapis.com/token",
              "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
              "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/serviceacct%40proj.iam.gserviceaccount.com",
              "universe_domain": "googleapis.com"
            }
            EOT
  resource_labels = {
    "somekey1"  = "v1"
    "otherkey1" = "v2"
  }
}
