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

resource "mezmo_pipeline" "my_pipeline" {
  title = "pipeline"
}
resource "mezmo_prometheus_remote_write_source" "metrics_source" {
  pipeline_id = mezmo_pipeline.my_pipeline.id
  title       = "My Prometheus Remote Write source"
  description = "This receives data from prometheus"
}
resource "mezmo_absence_alert" "no_data_alert_log_analysis" {
  pipeline_id             = mezmo_pipeline.my_pipeline.id
  component_kind          = "source"
  component_id            = mezmo_prometheus_remote_write_source.metrics_source.id
  inputs                  = [mezmo_prometheus_remote_write_source.metrics_source.id]
  name                    = "metrics absence alert"
  event_type              = "metric"
  window_duration_minutes = 15
  alert_payload = {
    service = {
      name          = "log_analysis"
      subject       = "No data received!"
      severity      = "WARNING"
      body          = "There has been no metrics data received in the last 15 minutes!"
      ingestion_key = "abc123"
    }
  }
}
resource "mezmo_absence_alert" "no_data_alert_webhook" {
  pipeline_id             = mezmo_pipeline.my_pipeline.id
  component_kind          = "source"
  component_id            = mezmo_prometheus_remote_write_source.metrics_source.id
  inputs                  = [mezmo_prometheus_remote_write_source.metrics_source.id]
  name                    = "metrics absence alert"
  event_type              = "metric"
  window_duration_minutes = 15
  alert_payload = {
    service = {
      name         = "webhook"
      uri          = "http://example.com/my_webhook"
      message_text = "There was an absence alert!"
      auth = {
        strategy = "bearer"
        token    = "abc123"
      }
      headers = {
        "x-my-header" = "header_value"
      }
    }
  }
}
