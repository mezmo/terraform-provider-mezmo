---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "mezmo_splunk_hec_logs_destination Resource - terraform-provider-mezmo"
subcategory: ""
description: |-
  Publishes log events to a Splunk HTTP Event Collector
---

# mezmo_splunk_hec_logs_destination (Resource)

Publishes log events to a Splunk HTTP Event Collector

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

variable "my_splunk_token" {
  type = string
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

resource "mezmo_splunk_hec_logs_destination" "destination1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My destination"
  description = "Send logs to a Splunk HEC server"
  inputs      = [mezmo_demo_source.source1.id]
  endpoint    = "https://example3.com"
  token       = var.my_splunk_token
  source = {
    value = "my source"
  }
  index = {
    field = ".my_index"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `endpoint` (String) The base URL for the Splunk instance. The collector path, such as `/services/collector/events`, will be automatically inferred from the destination's configuration.
- `pipeline_id` (String) The uuid of the pipeline
- `token` (String, Sensitive) The default token to authenticate to Splunk HEC

### Optional

- `ack_enabled` (Boolean) Acknowledge data from the source when it reaches the destination
- `compression` (String) The compression strategy used on the encoded data prior to sending
- `description` (String) A user-defined value describing the destination
- `host_field` (String) The field that contains the hostname to include in the event
- `index` (Attributes) The name of the index to send events to. Use the field path  "metadata.index" to use the upstream index value from a HEC log source (see [below for nested schema](#nestedatt--index))
- `inputs` (List of String) The ids of the input components
- `source` (Attributes) The source of events sent to this destination. This is typically the filename the logs originated from. Use the field path "metadata.source" to use the upstream source value from a HEC log source (see [below for nested schema](#nestedatt--source))
- `source_type` (Attributes) The sourcetype of events sent to this destination. Use the field path "metadata.sourcetype" to use the upstream sourcetype value from a HEC log source (see [below for nested schema](#nestedatt--source_type))
- `timestamp_field` (String) The field that contains the timestamp to include in the event
- `title` (String) A user-defined title for the destination
- `tls_verify_certificate` (Boolean) Verify TLS Certificate

### Read-Only

- `generation_id` (Number) An internal field used for component versioning
- `id` (String) The uuid of the destination

<a id="nestedatt--index"></a>
### Nested Schema for `index`

Optional:

- `field` (String) The field path to use
- `value` (String) The fixed value to use


<a id="nestedatt--source"></a>
### Nested Schema for `source`

Optional:

- `field` (String) The field path to use
- `value` (String) The fixed value to use


<a id="nestedatt--source_type"></a>
### Nested Schema for `source_type`

Optional:

- `field` (String) The field path to use
- `value` (String) The fixed value to use
