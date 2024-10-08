
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
resource "mezmo_webhook_source" "my_webhook" {
  pipeline_id = mezmo_pipeline.my_pipeline.id
  title       = "My webhook source"
  description = "This is a source populated via a webhook call"
}
resource "mezmo_threshold_alert" "order_count" {
  pipeline_id    = mezmo_pipeline.my_pipeline.id
  component_kind = "source"
  component_id   = mezmo_webhook_source.my_webhook.id
  inputs         = [mezmo_webhook_source.my_webhook.id]
  name           = "More orders than expected"
  event_type     = "log"
  operation      = "custom"
  script         = <<-EOSCRIPT
    function rollup(accum, event, metadata) {
      if (!accum.order_count) {
        accum.order_count = 0;
      }
      accum.order_count += event.num_ordered;
      return accum;
    }
    EOSCRIPT
  conditional = {
    expressions = [
      {
        field        = ".order_count"
        operator     = "greater"
        value_number = 100
      },
      {
        field        = ".status"
        operator     = "less_or_equal"
        value_number = 300
      }
    ],
  }
  window_type             = "tumbling"
  window_duration_minutes = 60
  alert_payload = {
    service = {
      name         = "pager_duty"
      uri          = "https://example.com/pager_duty_api"
      source       = "{{.my_source}}"
      routing_key  = "abc123"
      severity     = "CRITICAL"
      event_action = "trigger"
      summary      = "Check to make sure there are no errors in pricing, and no unexpected special offers were released."
    }
    throttling = {
      window_secs = 3600
      threshold   = 1
    }
  }
}
