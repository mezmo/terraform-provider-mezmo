---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "mezmo_unroll_processor Resource - terraform-provider-mezmo"
subcategory: ""
description: |-
  Takes an array of events and emits them all as individual events
---

# mezmo_unroll_processor (Resource)

Takes an array of events and emits them all as individual events

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

resource "mezmo_http_source" "curl" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My data stream"
  description = "Send Curl data to the pipeline point of entry URL"
  decoding    = "json"
}

resource "mezmo_unroll_processor" "processor1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My unroll processor"
  description = "I want events for each element of this array"
  inputs      = [mezmo_http_source.curl.id]
  field       = ".my_array_prop"
  values_only = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `field` (String) The field name that contains an array of events
- `pipeline_id` (String) The uuid of the pipeline

### Optional

- `description` (String) A user-defined value describing the processor
- `inputs` (List of String) The ids of the input components
- `title` (String) A user-defined title for the processor
- `values_only` (Boolean) When enabled, the values from the specified array field will be emitted as new events. Otherwise, the original event will be duplicated for each value in the array field, with the unrolled value present in the field specified.

### Read-Only

- `generation_id` (Number) An internal field used for component versioning
- `id` (String) The uuid of the processor
