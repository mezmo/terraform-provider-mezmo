package alerts

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
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
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
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
						Sensitive:   true,
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
						Sensitive:   true,
						Description: "The ingestion key for the service (Log Analysis).",
					},
					"auth": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "Configures HTTP authentication (Webhook).",
						Attributes: map[string]schema.Attribute{
							"strategy": schema.StringAttribute{
								Required:   true,
								Validators: []validator.String{stringvalidator.OneOf("basic", "bearer")},
							},
							"user": schema.StringAttribute{
								Optional:   true,
								Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"password": schema.StringAttribute{
								Sensitive:  true,
								Optional:   true,
								Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"token": schema.StringAttribute{
								Sensitive:  true,
								Optional:   true,
								Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
							},
						},
					},
					"headers": schema.MapAttribute{
						Optional:    true,
						Description: "Optional key/val request headers (Webhook).",
						ElementType: StringType{},
						Validators: []validator.Map{
							mapvalidator.All(
								mapvalidator.KeysAre(stringvalidator.LengthAtLeast(1)),
								mapvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
							),
						},
					},
					"method": schema.StringAttribute{
						Optional:    true,
						Description: "The HTTP method to use for the destination (Webhook, default is `post`).",
						Validators:  []validator.String{stringvalidator.OneOf("post", "put", "patch", "delete", "get", "head", "options", "trace")},
					},
				},
			},
			"throttling": schema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true, // There are defaults set by the server, hence Computed
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
		"auth": ObjectType{
			AttrTypes: map[string]attr.Type{
				"strategy": StringType{},
				"user":     StringType{},
				"password": StringType{},
				"token":    StringType{},
			},
		},
		"headers": MapType{
			ElemType: StringType{},
		},
		"method": StringType{},
	}
	serviceValues := map[string]attr.Value{
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
		"auth":          NewObjectNull(serviceTypes["auth"].(ObjectType).AttributeTypes()),
		"headers":       NewMapNull(StringType{}),
		"method":        NewStringNull(),
	}
	throttlingTypes := map[string]attr.Type{
		"window_secs": Int64Type{},
		"threshold":   Int64Type{},
	}
	throttlingValues := map[string]attr.Value{}
	for key, value := range component["service"].(map[string]any) {
		switch key {
		case "auth":
			auth, _ := value.(map[string]any)
			if auth["strategy"] != "none" {
				authTypes := serviceTypes["auth"].(ObjectType).AttributeTypes()
				authValues := MapAnyFillMissingValues(authTypes, auth, MapKeys(authTypes))
				serviceValues[key] = NewObjectValueMust(authTypes, authValues)
			}
		case "headers":
			headerArray, _ := value.([]any)
			if len(headerArray) > 0 {
				headerMap := make(map[string]any, len(headerArray))
				for _, obj := range headerArray {
					obj := obj.(map[string]any)
					key := obj["header_name"].(string)
					value := obj["header_value"].(string)
					headerMap[key] = value
				}
				serviceValues["headers"] = NewMapValueMust(serviceTypes["headers"].(MapType).ElementType(), MapAnyToMapValues(headerMap))
			}
		default:
			serviceValues[key] = NewStringValue(value.(string))
		}
	}
	for key, value := range component["throttling"].(map[string]any) {
		// All throttling values are float64/Int64
		throttlingValues[key] = NewInt64Value(int64(value.(float64)))
	}

	alertPayloadAttrs := NewObjectValueMust(map[string]attr.Type{
		"service":    ObjectType{AttrTypes: serviceTypes},
		"throttling": ObjectType{AttrTypes: throttlingTypes},
	}, map[string]attr.Value{
		"service":    NewObjectValueMust(serviceTypes, serviceValues),
		"throttling": NewObjectValueMust(throttlingTypes, throttlingValues),
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

	serviceValues := attrs["service"].(ObjectValue).Attributes()
	serviceName := serviceValues["name"].(StringValue).ValueString()
	service := map[string]any{
		"name": serviceName,
	}

	// Check that there are no extra properties that the service doesn't expect
	for attribute := range serviceValues {
		if attribute == "name" {
			continue
		}
		switch attribute {
		case "uri":
			if serviceName == "slack" || serviceName == "pager_duty" || serviceName == "webhook" {
				continue
			}
		case "message_text":
			if serviceName == "slack" || serviceName == "webhook" {
				continue
			}
		case "summary", "source", "routing_key", "event_action":
			if serviceName == "pager_duty" {
				continue
			}
		case "severity":
			if serviceName == "pager_duty" || serviceName == "log_analysis" {
				continue
			}
		case "subject", "body", "ingestion_key":
			if serviceName == "log_analysis" {
				continue
			}
		case "auth", "headers", "method":
			if serviceName == "webhook" {
				continue
			}
		}

		if serviceValues[attribute].IsNull() {
			// The TF framework will populate optional fields with null values. Values must
			// be non-null to be considered an error
			continue
		}
		dd.AddError(
			"Error in plan",
			fmt.Sprintf("Attribute `%s` is not allowed for service `%s`", attribute, serviceName),
		)
	}

	// Check for *missing* properties for the selected service, and construct the payload.
	// At this point, there are no extra properties which is validated above.
	switch serviceName {
	case "slack", "webhook":
		if serviceValues["uri"].(StringValue).IsNull() {
			dd.AddError(
				"Error in plan",
				fmt.Sprintf("`uri` is required for the `%s` service", serviceName),
			)
		}
		messageText := serviceValues["message_text"].(StringValue)
		if messageText.IsNull() {
			dd.AddError(
				"Error in plan",
				fmt.Sprintf("`message_text` is required for the `%s` service", serviceName),
			)
		}
		if authObj, ok := serviceValues["auth"]; ok && !authObj.IsNull() {
			auth := MapValuesToMapAny(authObj, dd)
			service["auth"] = auth
			if auth["strategy"] == "basic" {
				if auth["user"] == nil || auth["password"] == nil {
					dd.AddError(
						"Error in plan",
						"Basic auth requires user and password fields to be defined for the `webhook` service")
				}
			} else {
				if auth["token"] == nil {
					dd.AddError(
						"Error in plan",
						"Bearer auth requires token field to be defined for the `webhook` service")
				}
			}
		}
		if headersObj, ok := serviceValues["headers"]; ok && !headersObj.IsNull() {
			headerMap := MapValuesToMapAny(headersObj, dd)
			if len(headerMap) > 0 {
				headerArray := make([]map[string]string, 0, len(headerMap))
				for k, v := range headerMap {
					headerArray = append(headerArray, map[string]string{"header_name": k, "header_value": v.(string)})
				}
				service["headers"] = headerArray
			}
		}
		service["message_text"] = messageText.ValueString()
		service["uri"] = serviceValues["uri"].(StringValue).ValueString()
		method := serviceValues["method"].(StringValue)
		if !method.IsNull() {
			service["method"] = method.ValueString()
		}

	case "pager_duty":
		if serviceValues["uri"].(StringValue).IsNull() {
			dd.AddError(
				"Error in plan",
				fmt.Sprintf("`uri` is required for the `%s` service", serviceName),
			)
		}
		summary := serviceValues["summary"].(StringValue)
		severity := serviceValues["severity"].(StringValue)
		source := serviceValues["source"].(StringValue)
		routingKey := serviceValues["routing_key"].(StringValue)
		eventAction := serviceValues["event_action"].(StringValue)

		if summary.IsNull() {
			dd.AddError(
				"Error in plan",
				"`summary` is required for the `pager_duty` service",
			)
		}
		if severity.IsNull() {
			dd.AddError(
				"Error in plan",
				"`severity` is required for the `pager_duty` service",
			)
		}
		if severity.IsNull() {
			dd.AddError(
				"Error in plan",
				"`source` is required for the `pager_duty` service",
			)
		}
		if routingKey.IsNull() {
			dd.AddError(
				"Error in plan",
				"`routing_key` is required for the `pager_duty` service",
			)
		}
		if eventAction.IsNull() {
			dd.AddError(
				"Error in plan",
				"`event_action` is required for the `pager_duty` service",
			)
		}
		service["summary"] = summary.ValueString()
		service["severity"] = severity.ValueString()
		service["source"] = source.ValueString()
		service["routing_key"] = routingKey.ValueString()
		service["event_action"] = eventAction.ValueString()
		service["uri"] = serviceValues["uri"].(StringValue).ValueString()

	case "log_analysis":
		severity := serviceValues["severity"].(StringValue)
		subject := serviceValues["subject"].(StringValue)
		body := serviceValues["body"].(StringValue)
		ingestionKey := serviceValues["ingestion_key"].(StringValue)

		if severity.IsNull() {
			dd.AddError(
				"Error in plan",
				"`severity` is required for the `log_analysis` service",
			)
		}
		if subject.IsNull() {
			dd.AddError(
				"Error in plan",
				"`subject` is required for the `log_analysis` service",
			)
		}
		if body.IsNull() {
			dd.AddError(
				"Error in plan",
				"`body` is required for the `log_analysis` service",
			)
		}
		if ingestionKey.IsNull() {
			dd.AddError(
				"Error in plan",
				"`ingestion_key` is required for the `log_analysis` service",
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
		throttlingValues := throttlingObj.(ObjectValue).Attributes()
		throttling := map[string]any{}
		if windowSecsAttr, ok := throttlingValues["window_secs"]; ok && !windowSecsAttr.IsNull() {
			throttling["window_secs"] = windowSecsAttr.(Int64Value).ValueInt64()
		}
		if thresholdAttr, ok := throttlingValues["threshold"]; ok && !thresholdAttr.IsNull() {
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
