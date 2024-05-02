
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
      }
    ],
  }
  window_type             = "tumbling"
  window_duration_minutes = 60
  subject                 = "Lots of orders coming in the last hour!"
  severity                = "WARNING"
  message                 = "Check to make sure there are no errors in pricing, and no unexpected special offers were released."
  ingestion_key           = "abc123"
}
