package alerts

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	"severity": schema.StringAttribute{
		Optional:    true,
		Computed:    true,
		Description: "The severity level of the alert.",
		Validators: []validator.String{
			stringvalidator.OneOf("INFO", "WARNING", "ERROR", "CRITICAL"),
		},
		Default: stringdefault.StaticString("INFO"),
	},
	"style": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Description: "Configuration for how the alert message will be constructed. For " +
			"`static`, exact strings will be used. For `template`, the alert subjec and body " +
			"will allow for placeholders to substitute values from the event.",
		Validators: []validator.String{
			stringvalidator.OneOf("static", "template"),
		},
		Default: stringdefault.StaticString("static"),
	},
	"subject": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
			stringvalidator.LengthAtMost(200),
		},
		MarkdownDescription: "The subject line to use when the alert is sent. For a `template` style, " +
			"surround the field path in double curly braces.\n" +
			"```\n" + `{{"{{.my_field}} had a count of {{metadata.aggregate.event_count}}"}}` + "\n```",
	},
	"body": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
			stringvalidator.LengthAtMost(1024),
		},
		Description: "The message body to use when the alert is sent. For a `template` style, " +
			"surround the field path in double curly braces.\n" +
			"```\n" + `{{"{{.my_field}} had a count of {{metadata.aggregate.event_count}}"}}` + "\n```",
	},
	"ingestion_key": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
		Description: "The key required to ingest the alert into Log Analysis",
	},
}

type CheckedFields struct {
	Operation      StringValue
	EventType      StringValue
	Script         StringValue
	EventTimestamp StringValue
	GroupBy        ListValue
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
