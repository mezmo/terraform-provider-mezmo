---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "mezmo_datadog_metrics_destination Resource - terraform-provider-mezmo"
subcategory: ""
description: |-
  Publishes metric events to Datadog
---

# mezmo_datadog_metrics_destination (Resource)

Publishes metric events to Datadog

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

resource "mezmo_pipeline" "pipeline1" {
  title = "My pipeline"
}

resource "mezmo_demo_source" "source1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My source"
  description = "This is the point of entry for our data"
  format      = "json"
}

resource "mezmo_datadog_metrics_destination" "destination1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My destination"
  description = "Datadog metrics destination"
  site        = "us1"
  api_key     = "<secret-api-key>"
  inputs      = [mezmo_demo_source.source1.id]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `api_key` (String, Sensitive) Datadog metrics application API key.
- `pipeline_id` (String) The uuid of the pipeline
- `site` (String) The Datadog site (region) to send metrics to.

### Optional

- `ack_enabled` (Boolean) Acknowledge data from the source when it reaches the destination
- `description` (String) A user-defined value describing the destination
- `inputs` (List of String) The ids of the input components
- `title` (String) A user-defined title for the destination

### Read-Only

- `generation_id` (Number) An internal field used for component versioning
- `id` (String) The uuid of the destination
