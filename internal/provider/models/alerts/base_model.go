package alerts

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type SchemaAttributes map[string]schema.Attribute

var baseAlertSchemaAttributes = SchemaAttributes{
	// Non-config fields
	"id": schema.StringAttribute{
		Computed:    true,
		Description: "The uuid of the alert",
	},
	"pipeline_id": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
		Description: "The uuid of the pipeline",
	},
	"component_kind": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("source", "transform"),
		},
		Description: "The kind of component that the alert is attached to",
	},
	"component_id": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
		Description: "The uuid of the component that the alert is attached to",
	},
	"inputs": schema.ListAttribute{
		ElementType: types.StringType,
		Required:    true,
		Validators: []validator.List{
			listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
			listvalidator.SizeAtLeast(1),
			listvalidator.UniqueValues(),
		},
		Description: "The ids of the input components. This could be the id of a match arm " +
			"for a route processor, or simply the id of the component.",
	},
	"active": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Indicates if the alert is turned on or off",
		Default:     booldefault.StaticBool(true),
	},
	// Alert config fields
	"name": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
			stringvalidator.LengthAtMost(200),
		},
		Description: "The name of the alert.",
	},
	"description": schema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			stringvalidator.LengthAtMost(1024),
		},
		Description: "An optional description describing what the alert is for.",
	},
	"event_type": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("log", "metric"),
		},
		Description: "The type of event is either a Log event or a Metric event.",
	},
	"group_by": schema.ListAttribute{
		ElementType: StringType{},
		Optional:    true,
		Validators: []validator.List{
			listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
			listvalidator.ValueStringsAre(stringvalidator.LengthAtMost(100)),
			listvalidator.SizeAtMost(50),
			listvalidator.UniqueValues(),
		},
		Description: "When aggregating, group events based on matching values from each of " +
			"these field paths. Supports nesting via dot-notation. This value is " +
			"optional for Metric event types, and SHOULD be used for Log event types.",
	},
	"window_type": schema.StringAttribute{
		Optional:   true,
		Computed:   true,
		Validators: []validator.String{stringvalidator.OneOf("tumbling", "sliding")},
		Default:    stringdefault.StaticString("tumbling"),
		Description: "Sliding windows can overlap, whereas tumbling windows are disjoint. " +
			"For example, a tumbling window has a fixed time span and any events that " +
			"fall within the \"window duration\" will be used in the aggregate. " +
			"In a sliding window, the aggregation occurs every \"window duration\" seconds " +
			"after an event is encountered.",
	},
	"window_duration_minutes": schema.Int64Attribute{
		Optional: true,
		Computed: true,
		Default:  int64default.StaticInt64(5),
		Validators: []validator.Int64{
			int64validator.AtLeast(1),
			int64validator.AtMost(1440),
		},
		Description: "The duration of the aggregation window in minutes.",
	},
	"event_timestamp": schema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
			stringvalidator.LengthAtMost(100),
		},
		Description: "The path to a field on the event that contains an epoch timestamp " +
			"value. If an event does not have a timestamp field, events will be associated " +
			"to the wall clock value when the event is processed. " +
			"Required for Log event types and disallowed for Metric event types.",
	},
	"alert_payload": schema.SingleNestedAttribute{
		Required: true,
		Description: "Configure where the alert will be sent, including choosing a service " +
			"and throttling options. All options for the chosen `service` will be required. " +
			"All text fields support templating.",
		Attributes: map[string]schema.Attribute{
			"service": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Configuration for the service receiving the alert.",
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required:    true,
						Description: "The name of the service.",
						Validators: []validator.String{
							stringvalidator.OneOf("slack", "pager_duty", "webhook", "log_analysis"),
						},
					},
					"uri": schema.StringAttribute{
						Optional:    true,
						Description: "The URI of the service (Slack, PagerDuty, Webhook).",
					},
					"message_text": schema.StringAttribute{
						Optional:    true,
						Description: "The text value of the notification message (Slack, Webhook).",
					},
					"summary": schema.StringAttribute{
						Optional:    true,
						Description: "Summarize the alert details (PagerDuty).",
					},
					"source": schema.StringAttribute{
						Optional:    true,
						Description: "The source of the alert (PagerDuty).",
					},
					"routing_key": schema.StringAttribute{
						Optional:    true,
						Description: "The service's routing key (PagerDuty).",
					},
					"event_action": schema.StringAttribute{
						Optional:    true,
						Description: "The event action to use (PagerDuty).",
					},
					"severity": schema.StringAttribute{
						Optional:    true,
						Description: "The severity level of the alert (PagerDuty, Log Analysis).",
						Validators: []validator.String{
							stringvalidator.OneOf("INFO", "WARNING", "ERROR", "CRITICAL"),
						},
					},
					"subject": schema.StringAttribute{
						Optional:    true,
						Description: "The main subject line of the message (Log Analysis).",
					},
					"body": schema.StringAttribute{
						Optional:    true,
						Description: "Additional information to be added to the message (Log Analysis).",
					},
					"ingestion_key": schema.StringAttribute{
						Optional:    true,
						Description: "The ingestion key for the service (Log Analysis).",
					},
				},
			},
			"throttling": schema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Configure throttling options for the service receiving the alert.",
				Attributes: map[string]schema.Attribute{
					"window_secs": schema.Int64Attribute{
						Optional: true,
						Computed: true,
						Description: "The time frame during which the number of notifications " +
							"(set by the `threshold`) is permitted (default: 60).",
						Validators: []validator.Int64{
							int64validator.AtLeast(1),
						},
					},
					"threshold": schema.Int64Attribute{
						Optional: true,
						Computed: true,
						Description: "The maximum number of notifications allowed over the given time " +
							"window set by `window_secs` (default: 1).",
						Validators: []validator.Int64{
							int64validator.AtLeast(1),
						},
					},
				},
			},
		},
	},
}

type CheckedFields struct {
	Operation      StringValue
	EventType      StringValue
	Script         StringValue
	EventTimestamp StringValue
	GroupBy        ListValue
}

// Convert the api response for `alert_payload` into a Terraform model
func GetAlertPayloadToModel(component map[string]any) ObjectValue {
	// All properties need to be defined regardless of service name. Initialize with null values.
	serviceTypes := map[string]attr.Type{
		"name":          StringType{},
		"uri":           StringType{},
		"message_text":  StringType{},
		"summary":       StringType{},
		"source":        StringType{},
		"routing_key":   StringType{},
		"event_action":  StringType{},
		"severity":      StringType{},
		"subject":       StringType{},
		"body":          StringType{},
		"ingestion_key": StringType{},
	}
	serviceAttrs := map[string]attr.Value{
		"name":          NewStringNull(),
		"uri":           NewStringNull(),
		"message_text":  NewStringNull(),
		"summary":       NewStringNull(),
		"source":        NewStringNull(),
		"routing_key":   NewStringNull(),
		"event_action":  NewStringNull(),
		"severity":      NewStringNull(),
		"subject":       NewStringNull(),
		"body":          NewStringNull(),
		"ingestion_key": NewStringNull(),
	}
	throttlingTypes := map[string]attr.Type{
		"window_secs": Int64Type{},
		"threshold":   Int64Type{},
	}
	throttlingAttrs := map[string]attr.Value{}
	for key, value := range component["service"].(map[string]any) {
		// All service values are strings
		serviceAttrs[key] = NewStringValue(value.(string))
	}
	for key, value := range component["throttling"].(map[string]any) {
		// All throttling values are float64/Int64
		throttlingAttrs[key] = NewInt64Value(int64(value.(float64)))
	}

	alertPayloadAttrs := NewObjectValueMust(map[string]attr.Type{
		"service":    ObjectType{AttrTypes: serviceTypes},
		"throttling": ObjectType{AttrTypes: throttlingTypes},
	}, map[string]attr.Value{
		"service":    NewObjectValueMust(serviceTypes, serviceAttrs),
		"throttling": NewObjectValueMust(throttlingTypes, throttlingAttrs),
	})

	return alertPayloadAttrs
}

// Assemble the `alert_payload` object to be sent to the service, including
// error checking. The payload must match the required structure for the chosen service,
// i.e. there can be no extra fields for options that the service `name` doesn't expect.
func GetAlertPayloadFromModel(v attr.Value, dd *diag.Diagnostics) map[string]any {
	value, ok := v.(ObjectValue)
	if !ok {
		panic(fmt.Errorf("Expected an object but did not receive one: %+v", v))
	}
	attrs := value.Attributes()

	serviceAttrs := attrs["service"].(ObjectValue).Attributes()
	serviceName := serviceAttrs["name"].(StringValue).ValueString()
	service := map[string]any{
		"name": serviceName,
	}

	if serviceName != "log_analysis" && serviceAttrs["uri"].(StringValue).IsNull() {
		dd.AddError(
			"Error in plan",
			"`uri` is required for Slack, PagerDuty, or Webhook",
		)
	}

	switch serviceName {
	case "slack", "webhook":
		messageText := serviceAttrs["message_text"].(StringValue)
		if messageText.IsNull() {
			dd.AddError(
				"Error in plan",
				"`message_text` is required for Slack or Webhook notifications",
			)
		}
		service["message_text"] = messageText.ValueString()
		service["uri"] = serviceAttrs["uri"].(StringValue).ValueString()

	case "pager_duty":
		summary := serviceAttrs["summary"].(StringValue)
		severity := serviceAttrs["severity"].(StringValue)
		source := serviceAttrs["source"].(StringValue)
		routingKey := serviceAttrs["routing_key"].(StringValue)
		eventAction := serviceAttrs["event_action"].(StringValue)

		if summary.IsNull() {
			dd.AddError(
				"Error in plan",
				"`summary` is required for PagerDuty notifications",
			)
		}
		if severity.IsNull() {
			dd.AddError(
				"Error in plan",
				"`severity` is required for PagerDuty notifications",
			)
		}
		if severity.IsNull() {
			dd.AddError(
				"Error in plan",
				"`source` is required for PagerDuty notifications",
			)
		}
		if routingKey.IsNull() {
			dd.AddError(
				"Error in plan",
				"`routing_key` is required for PagerDuty notifications",
			)
		}
		if eventAction.IsNull() {
			dd.AddError(
				"Error in plan",
				"`event_action` is required for PagerDuty notifications",
			)
		}
		service["summary"] = summary.ValueString()
		service["severity"] = severity.ValueString()
		service["source"] = source.ValueString()
		service["routing_key"] = routingKey.ValueString()
		service["event_action"] = eventAction.ValueString()
		service["uri"] = serviceAttrs["uri"].(StringValue).ValueString()

	case "log_analysis":
		severity := serviceAttrs["severity"].(StringValue)
		subject := serviceAttrs["subject"].(StringValue)
		body := serviceAttrs["body"].(StringValue)
		ingestionKey := serviceAttrs["ingestion_key"].(StringValue)

		if severity.IsNull() {
			dd.AddError(
				"Error in plan",
				"`severity` is required for Log Analysis notifications",
			)
		}
		if subject.IsNull() {
			dd.AddError(
				"Error in plan",
				"`subject` is required for Log Analysis notifications",
			)
		}
		if body.IsNull() {
			dd.AddError(
				"Error in plan",
				"`body` is required for Log Analysis notifications",
			)
		}
		if ingestionKey.IsNull() {
			dd.AddError(
				"Error in plan",
				"`ingestion_key` is required for Log Analysis notifications",
			)
		}
		service["severity"] = severity.ValueString()
		service["subject"] = subject.ValueString()
		service["body"] = body.ValueString()
		service["ingestion_key"] = ingestionKey.ValueString()
	}

	// Create the full payload, including `service` and `throttling` (if provided).
	// `null` values cannot be sent, so we must be surgical about the object's structure.
	alertPayload := map[string]any{
		"service": service,
	}

	if throttlingObj, ok := attrs["throttling"]; ok && !throttlingObj.IsNull() {
		throttlingAttrs := throttlingObj.(ObjectValue).Attributes()
		throttling := map[string]any{}
		if windowSecsAttr, ok := throttlingAttrs["window_secs"]; ok && !windowSecsAttr.IsNull() {
			throttling["window_secs"] = windowSecsAttr.(Int64Value).ValueInt64()
		}
		if thresholdAttr, ok := throttlingAttrs["threshold"]; ok && !thresholdAttr.IsNull() {
			throttling["threshold"] = thresholdAttr.(Int64Value).ValueInt64()
		}
		alertPayload["throttling"] = throttling

	}

	return alertPayload
}

func OperationAndScriptErrorChecks(plan *CheckedFields, dd *diag.Diagnostics) *diag.Diagnostics {
	operation := plan.Operation.ValueString()
	eventType := plan.EventType.ValueString()

	// Errors agnostic of event type
	if operation == "custom" {
		if plan.Script.IsNull() {
			dd.AddError(
				"Error in plan",
				"A 'custom' operation requires a valid JS `script` function",
			)
		}
	} else if !plan.Script.IsNull() {
		dd.AddError(
			"Error in plan",
			"`script` cannot be set when `operation` is not 'custom'",
		)
	}

	if eventType == "log" {
		if operation != "custom" {
			dd.AddError(
				"Error in plan",
				"A 'log' event type requires a 'custom' `operation` and a valid JS `script` function",
			)
		}
	}
	return dd
}

func CustomErrorChecks(plan *CheckedFields, dd *diag.Diagnostics) *diag.Diagnostics {
	eventType := plan.EventType.ValueString()

	if eventType == "log" {
		if plan.EventTimestamp.IsNull() {
			dd.AddWarning(
				"A 'log' event should specify an `event_timestamp` field",
				"If the event payload does not have an epoch timestamp field, then the "+
					"system time will be used at the time the event is processed.",
			)
		}
		if plan.GroupBy.IsNull() {
			dd.AddWarning(
				"A 'log' event should use a `group_by` field",
				"If a `group_by` is not specified, then all events will be combined, elimiating "+
					"their distinctness. The recommendation is [\".name\", \".namespace\", \".tags\"]",
			)
		}
	} else {
		// metric
		if !plan.EventTimestamp.IsNull() {
			dd.AddError(
				"Error in plan",
				"A 'metric' event type cannot have an `event_timestamp` field",
			)
		}
	}
	return dd
}

func ExtendBaseAttributes(target SchemaAttributes) SchemaAttributes {
	return ExtendSchemaAttributes(baseAlertSchemaAttributes, target)
}

func ExtendSchemaAttributes(fromAttributes SchemaAttributes, toAttributes SchemaAttributes) SchemaAttributes {
	for k, v := range fromAttributes {
		toAttributes[k] = v
	}
	return toAttributes
}
