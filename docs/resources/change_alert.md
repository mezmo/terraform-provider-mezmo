---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "mezmo_change_alert Resource - terraform-provider-mezmo"
subcategory: ""
description: |-
  Represents a Change Alert in a Pipeline
---

# mezmo_change_alert (Resource)

Represents a Change Alert in a Pipeline

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
  subject                 = "Spike in ordering!"
  severity                = "WARNING"
  body                    = "There has been a > 20% increase in orders over the last 15 minutes. Check application scaling."
  ingestion_key           = "abc123"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `body` (String) The message body to use when the alert is sent. For a `template` style, surround the field path in double curly braces.
```
{{.my_field}} had a count of {{metadata.aggregate.event_count}}
```
- `component_id` (String) The uuid of the component that the alert is attached to
- `component_kind` (String) The kind of component that the alert is attached to
- `conditional` (Attributes) A group of expressions (optionally nested) joined by a logical operator (see [below for nested schema](#nestedatt--conditional))
- `event_type` (String) The type of event is either a Log event or a Metric event.
- `ingestion_key` (String) The key required to ingest the alert into Log Analysis
- `inputs` (List of String) The ids of the input components. This could be the id of a match arm for a route processor, or simply the id of the component.
- `name` (String) The name of the alert.
- `operation` (String) Specifies the type of aggregation operation to use with the window type and duration. This value must be `custom` for a Log event type.
- `pipeline_id` (String) The uuid of the pipeline
- `subject` (String) The subject line to use when the alert is sent. For a `template` style, surround the field path in double curly braces.
```
{{.my_field}} had a count of {{metadata.aggregate.event_count}}
```

### Optional

- `active` (Boolean) Indicates if the alert is turned on or off
- `description` (String) An optional description describing what the alert is for.
- `event_timestamp` (String) The path to a field on the event that contains an epoch timestamp value. If an event does not have a timestamp field, events will be associated to the wall clock value when the event is processed. Required for Log event types and disallowed for Metric event types.
- `group_by` (List of String) When aggregating, group events based on matching values from each of these field paths. Supports nesting via dot-notation. This value is optional for Metric event types, and SHOULD be used for Log event types.
- `script` (String) A custom JavaScript function that will control the aggregation. At the time of flushing, this aggregation will become the emitted event. This script is required when choosing a `custom` operation.
- `severity` (String) The severity level of the alert.
- `style` (String) Configuration for how the alert message will be constructed. For `static`, exact strings will be used. For `template`, the alert subjec and body will allow for placeholders to substitute values from the event.
- `window_duration_minutes` (Number) The duration of the aggregation window in minutes.
- `window_type` (String) Sliding windows can overlap, whereas tumbling windows are disjoint. For example, a tumbling window has a fixed time span and any events that fall within the "window duration" will be used in the aggregate. In a sliding window, the aggregation occurs every "window duration" seconds after an event is encountered.

### Read-Only

- `id` (String) The uuid of the alert

<a id="nestedatt--conditional"></a>
### Nested Schema for `conditional`

Optional:

- `expressions` (Attributes List) Defines a list of expressions for field comparisons (see [below for nested schema](#nestedatt--conditional--expressions))
- `expressions_group` (Attributes List) A group of expressions joined by a logical operator (see [below for nested schema](#nestedatt--conditional--expressions_group))
- `logical_operation` (String) The logical operation (AND/OR) to be applied to the list of conditionals

<a id="nestedatt--conditional--expressions"></a>
### Nested Schema for `conditional.expressions`

Required:

- `field` (String) The field path whose value will be used in the comparison
- `operator` (String) The comparison operator. Possible values are: percent_change_greater, percent_change_greater_or_equal, percent_change_less, percent_change_less_or_equal, value_change_greater, value_change_greater_or_equal, value_change_less or value_change_less_or_equal.

Optional:

- `value_number` (Number) The operand to compare the field value with, when the value is a number
- `value_string` (String) The operand to compare the field value with, when the value is a string


<a id="nestedatt--conditional--expressions_group"></a>
### Nested Schema for `conditional.expressions_group`

Optional:

- `expressions` (Attributes List) Defines a list of expressions for field comparisons (see [below for nested schema](#nestedatt--conditional--expressions_group--expressions))
- `expressions_group` (Attributes List) A group of expressions joined by a logical operator (see [below for nested schema](#nestedatt--conditional--expressions_group--expressions_group))
- `logical_operation` (String) The logical operation (AND/OR) to be applied to the list of conditionals

<a id="nestedatt--conditional--expressions_group--expressions"></a>
### Nested Schema for `conditional.expressions_group.expressions`

Required:

- `field` (String) The field path whose value will be used in the comparison
- `operator` (String) The comparison operator. Possible values are: percent_change_greater, percent_change_greater_or_equal, percent_change_less, percent_change_less_or_equal, value_change_greater, value_change_greater_or_equal, value_change_less or value_change_less_or_equal.

Optional:

- `value_number` (Number) The operand to compare the field value with, when the value is a number
- `value_string` (String) The operand to compare the field value with, when the value is a string


<a id="nestedatt--conditional--expressions_group--expressions_group"></a>
### Nested Schema for `conditional.expressions_group.expressions_group`

Optional:

- `expressions` (Attributes List) Defines a list of expressions for field comparisons (see [below for nested schema](#nestedatt--conditional--expressions_group--expressions_group--expressions))
- `expressions_group` (Attributes List) A group of expressions joined by a logical operator (see [below for nested schema](#nestedatt--conditional--expressions_group--expressions_group--expressions_group))
- `logical_operation` (String) The logical operation (AND/OR) to be applied to the list of conditionals

<a id="nestedatt--conditional--expressions_group--expressions_group--expressions"></a>
### Nested Schema for `conditional.expressions_group.expressions_group.logical_operation`

Required:

- `field` (String) The field path whose value will be used in the comparison
- `operator` (String) The comparison operator. Possible values are: percent_change_greater, percent_change_greater_or_equal, percent_change_less, percent_change_less_or_equal, value_change_greater, value_change_greater_or_equal, value_change_less or value_change_less_or_equal.

Optional:

- `value_number` (Number) The operand to compare the field value with, when the value is a number
- `value_string` (String) The operand to compare the field value with, when the value is a string


<a id="nestedatt--conditional--expressions_group--expressions_group--expressions_group"></a>
### Nested Schema for `conditional.expressions_group.expressions_group.logical_operation`

Optional:

- `expressions` (Attributes List) Defines a list of expressions for field comparisons (see [below for nested schema](#nestedatt--conditional--expressions_group--expressions_group--logical_operation--expressions))
- `expressions_group` (Attributes List) A group of expressions joined by a logical operator (see [below for nested schema](#nestedatt--conditional--expressions_group--expressions_group--logical_operation--expressions_group))
- `logical_operation` (String) The logical operation (AND/OR) to be applied to the list of conditionals

<a id="nestedatt--conditional--expressions_group--expressions_group--logical_operation--expressions"></a>
### Nested Schema for `conditional.expressions_group.expressions_group.logical_operation.expressions`

Required:

- `field` (String) The field path whose value will be used in the comparison
- `operator` (String) The comparison operator. Possible values are: percent_change_greater, percent_change_greater_or_equal, percent_change_less, percent_change_less_or_equal, value_change_greater, value_change_greater_or_equal, value_change_less or value_change_less_or_equal.

Optional:

- `value_number` (Number) The operand to compare the field value with, when the value is a number
- `value_string` (String) The operand to compare the field value with, when the value is a string


<a id="nestedatt--conditional--expressions_group--expressions_group--logical_operation--expressions_group"></a>
### Nested Schema for `conditional.expressions_group.expressions_group.logical_operation.expressions_group`

Optional:

- `expressions` (Attributes List) Defines a list of expressions for field comparisons (see [below for nested schema](#nestedatt--conditional--expressions_group--expressions_group--logical_operation--expressions_group--expressions))
- `expressions_group` (Attributes List) A group of expressions joined by a logical operator (see [below for nested schema](#nestedatt--conditional--expressions_group--expressions_group--logical_operation--expressions_group--expressions_group))
- `logical_operation` (String) The logical operation (AND/OR) to be applied to the list of conditionals

<a id="nestedatt--conditional--expressions_group--expressions_group--logical_operation--expressions_group--expressions"></a>
### Nested Schema for `conditional.expressions_group.expressions_group.logical_operation.expressions_group.logical_operation`

Required:

- `field` (String) The field path whose value will be used in the comparison
- `operator` (String) The comparison operator. Possible values are: percent_change_greater, percent_change_greater_or_equal, percent_change_less, percent_change_less_or_equal, value_change_greater, value_change_greater_or_equal, value_change_less or value_change_less_or_equal.

Optional:

- `value_number` (Number) The operand to compare the field value with, when the value is a number
- `value_string` (String) The operand to compare the field value with, when the value is a string


<a id="nestedatt--conditional--expressions_group--expressions_group--logical_operation--expressions_group--expressions_group"></a>
### Nested Schema for `conditional.expressions_group.expressions_group.logical_operation.expressions_group.logical_operation`

Optional:

- `expressions` (Attributes List) Defines a list of expressions for field comparisons (see [below for nested schema](#nestedatt--conditional--expressions_group--expressions_group--logical_operation--expressions_group--logical_operation--expressions))
- `expressions_group` (Attributes List) A group of expressions joined by a logical operator (see [below for nested schema](#nestedatt--conditional--expressions_group--expressions_group--logical_operation--expressions_group--logical_operation--expressions_group))
- `logical_operation` (String) The logical operation (AND/OR) to be applied to the list of conditionals

<a id="nestedatt--conditional--expressions_group--expressions_group--logical_operation--expressions_group--logical_operation--expressions"></a>
### Nested Schema for `conditional.expressions_group.expressions_group.logical_operation.expressions_group.logical_operation.logical_operation`

Required:

- `field` (String) The field path whose value will be used in the comparison
- `operator` (String) The comparison operator. Possible values are: percent_change_greater, percent_change_greater_or_equal, percent_change_less, percent_change_less_or_equal, value_change_greater, value_change_greater_or_equal, value_change_less or value_change_less_or_equal.

Optional:

- `value_number` (Number) The operand to compare the field value with, when the value is a number
- `value_string` (String) The operand to compare the field value with, when the value is a string


<a id="nestedatt--conditional--expressions_group--expressions_group--logical_operation--expressions_group--logical_operation--expressions_group"></a>
### Nested Schema for `conditional.expressions_group.expressions_group.logical_operation.expressions_group.logical_operation.logical_operation`

Optional:

- `expressions` (Attributes List) Defines a list of expressions for field comparisons (see [below for nested schema](#nestedatt--conditional--expressions_group--expressions_group--logical_operation--expressions_group--logical_operation--logical_operation--expressions))
- `logical_operation` (String) The logical operation (AND/OR) to be applied to the list of conditionals

<a id="nestedatt--conditional--expressions_group--expressions_group--logical_operation--expressions_group--logical_operation--logical_operation--expressions"></a>
### Nested Schema for `conditional.expressions_group.expressions_group.logical_operation.expressions_group.logical_operation.logical_operation.logical_operation`

Required:

- `field` (String) The field path whose value will be used in the comparison
- `operator` (String) The comparison operator. Possible values are: percent_change_greater, percent_change_greater_or_equal, percent_change_less, percent_change_less_or_equal, value_change_greater, value_change_greater_or_equal, value_change_less or value_change_less_or_equal.

Optional:

- `value_number` (Number) The operand to compare the field value with, when the value is a number
- `value_string` (String) The operand to compare the field value with, when the value is a string
