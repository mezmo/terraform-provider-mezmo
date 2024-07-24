
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
resource "mezmo_change_alert" "order_spike" {
  pipeline_id    = mezmo_pipeline.my_pipeline.id
  component_kind = "source"
  component_id   = mezmo_webhook_source.my_webhook.id
  inputs         = [mezmo_webhook_source.my_webhook.id]
  name           = "Spike in orders"
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
        operator     = "percent_change_greater"
        value_number = 20
      }
    ],
  }
  window_type             = "sliding"
  window_duration_minutes = 15

  alert_payload = {
    service = {
      name         = "webhook"
      uri          = "https://example.com/my_webhook"
      message_text = "There has been a > 20% increase in orders ({{.total_orders}}) over the last 15 minutes. Check application scaling."
    }
    throttling = {
      window_secs = 3600
      threshold   = 2
    }
  }
}
