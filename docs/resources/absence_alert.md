---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "mezmo_absence_alert Resource - terraform-provider-mezmo"
subcategory: ""
description: |-
  Represents an Absence Alert in a Pipeline
---

# mezmo_absence_alert (Resource)

Represents an Absence Alert in a Pipeline

## Example Usage

```terraform
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
resource "mezmo_absence_alert" "no_data_alert" {
  pipeline_id             = mezmo_pipeline.my_pipeline.id
  component_kind          = "source"
  component_id            = mezmo_prometheus_remote_write_source.metrics_source.id
  inputs                  = [mezmo_prometheus_remote_write_source.metrics_source.id]
  name                    = "metrics absence alert"
  event_type              = "metric"
  operation               = "sum"
  window_duration_minutes = 15
  subject                 = "No data received!"
  severity                = "WARNING"
  message                 = "There has been no metrics data recieved in the last 15 minutes!"
  ingestion_key           = "abc123"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `component_id` (String) The uuid of the component that the alert is attached to
- `component_kind` (String) The kind of component that the alert is attached to
- `event_type` (String) The type of event is either a Log event or a Metric event.
- `ingestion_key` (String) The key required to ingest the alert into Log Analysis
- `inputs` (List of String) The ids of the input components. This could be the id of a match arm for a route processor, or simply the id of the component.
- `message` (String) The message body to use when the alert is sent.
- `name` (String) The name of the alert.
- `operation` (String) Specifies the type of aggregation operation to use with the window type and duration. This value must be `custom` for a Log event type.
- `pipeline_id` (String) The uuid of the pipeline
- `subject` (String) The subject line to use when the alert is sent.

### Optional

- `active` (Boolean) Indicates if the alert is turned on or off
- `description` (String) An optional description describing what the alert is for.
- `event_timestamp` (String) The path to a field on the event that contains an epoch timestamp value. If an event does not have a timestamp field, events will be associated to the wall clock value when the event is processed. Required for Log event types and disallowed for Metric event types.
- `group_by` (List of String) When aggregating, group events based on matching values from each of these field paths. Supports nesting via dot-notation. This value is optional for Metric event types, and SHOULD be used for Log event types.
- `script` (String) A custom JavaScript function that will control the aggregation. At the time of flushing, this aggregation will become the emitted event. This script is required when choosing a `custom` operation.
- `severity` (String) The severity level of the alert.
- `style` (String) Configuration for how the alert message will be constructed.
- `window_duration_minutes` (Number) The duration of the aggregation window in minutes.
- `window_type` (String) Sliding windows can overlap, whereas tumbling windows are disjoint. For example, a tumbling window has a fixed time span and any events that fall within the "window duration" will be used in the aggregate. In a sliding window, the aggregation occurs every "window duration" seconds after an event is encountered.

### Read-Only

- `id` (String) The uuid of the alert