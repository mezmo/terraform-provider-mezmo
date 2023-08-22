---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "mezmo_splunk_hec_source Resource - terraform-provider-mezmo"
subcategory: ""
description: |-
  Receive Splunk logs
---

# mezmo_splunk_hec_source (Resource)

Receive Splunk logs



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `pipeline_id` (String) The uuid of the pipeline

### Optional

- `capture_metadata` (Boolean) Enable the inclusion of all http headers and query string parameters that were sent from the source
- `description` (String) A user-defined value describing the source component
- `gateway_route_id` (String) The uuid of a pre-existing source to be used as the input for this component. This can only be provided on resource creation (not update).
- `title` (String) A user-defined title for the source component

### Read-Only

- `generation_id` (Number) An internal field used for component versioning
- `id` (String) The uuid of the source component

