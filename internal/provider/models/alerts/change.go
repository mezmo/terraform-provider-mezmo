package alerts

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/models/modelutils"
)

const ALERT_TYPE_CHANGE = "change"

type ChangeAlertModel struct {
	Id                StringValue `tfsdk:"id"`
	PipelineId        StringValue `tfsdk:"pipeline_id"`
	ComponentKind     StringValue `tfsdk:"component_kind"`
	ComponentId       StringValue `tfsdk:"component_id"`
	Inputs            ListValue   `tfsdk:"inputs"`
	Active            BoolValue   `tfsdk:"active"`
	Name              StringValue `tfsdk:"name" user_config:"true"`
	Description       StringValue `tfsdk:"description" user_config:"true"`
	EventType         StringValue `tfsdk:"event_type"`
	GroupBy           ListValue   `tfsdk:"group_by" user_config:"true"`
	Operation         StringValue `tfsdk:"operation" user_config:"true"`
	Conditional       ObjectValue `tfsdk:"conditional" user_config:"true"`
	WindowType        StringValue `tfsdk:"window_type" user_config:"true"`
	WindowDurationMin Int64Value  `tfsdk:"window_duration_minutes" user_config:"true"`
	Script            StringValue `tfsdk:"script" user_config:"true"`
	EventTimestamp    StringValue `tfsdk:"event_timestamp" user_config:"true"`
	AlertPayload      ObjectValue `tfsdk:"alert_payload" user_config:"true"`
}

var ChangeAlertResourceSchema = schema.Schema{
	Description: "Represents a Change Alert in a Pipeline",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"conditional": schema.SingleNestedAttribute{
			Required:    true, // ParentConditionalAttribute is not required by default
			Description: ParentConditionalAttribute(Change_Operator_Labels).Description,
			Attributes:  ParentConditionalAttribute(Change_Operator_Labels).Attributes,
		},
		"operation": schema.StringAttribute{
			Required: true,
			Description: "Specifies the type of aggregation operation to use with the window type and duration. " +
				"This value must be `custom` for a Log event type.",
			Validators: []validator.String{
				stringvalidator.OneOf(MapKeys(Aggregate_Operations)...),
			},
		},
		"script": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.LengthAtMost(5000),
			},
			Description: "A custom JavaScript function that will control the aggregation. At the " +
				"time of flushing, this aggregation will become the emitted event. This script " +
				"is required when choosing a `custom` operation.",
		},
	}),
}

// From terraform schema/model to a struct for sending to the API
func ChangeAlertFromModel(plan *ChangeAlertModel, previousState *ChangeAlertModel) (*Alert, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	alertPayload := GetAlertPayloadFromModel(plan.AlertPayload, &dd)

	checkFields := CheckedFields{
		Operation:      plan.Operation,
		EventType:      plan.EventType,
		Script:         plan.Script,
		EventTimestamp: plan.EventTimestamp,
		GroupBy:        plan.GroupBy,
	}
	CustomErrorChecks(&checkFields, &dd)
	OperationAndScriptErrorChecks(&checkFields, &dd)

	if dd.HasError() {
		return nil, dd
	}

	// Inputs are required, so we can assemble them first
	inputs := make([]string, 0)
	for _, v := range plan.Inputs.Elements() {
		value, _ := v.(StringValue)
		inputs = append(inputs, value.ValueString())
	}

	component := Alert{
		PipelineId:    plan.PipelineId.ValueString(),
		ComponentKind: plan.ComponentKind.ValueString(),
		ComponentId:   plan.ComponentId.ValueString(),
		Inputs:        inputs,
		AlertConfig: map[string]any{
			"general": map[string]any{
				"name": plan.Name.ValueString(),
			},
			"evaluation": map[string]any{
				"alert_type":  ALERT_TYPE_CHANGE, // Required for the API, but hidden from the user here
				"event_type":  plan.EventType.ValueString(),
				"operation":   Aggregate_Operations[plan.Operation.ValueString()],
				"conditional": UnwindConditionalFromModel(plan.Conditional),
			},
			"alert_payload": alertPayload,
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
	}

	// Create all optional properties if they exist. The values will be:
	// .isUnknown() if it's computed and optional
	// .isNull() if it's optional and not set

	if !plan.Active.IsUnknown() {
		component.Active = plan.Active.ValueBool()
	}

	general := component.AlertConfig["general"].(map[string]any)
	evaluation := component.AlertConfig["evaluation"].(map[string]any)

	if !plan.Description.IsNull() {
		general["description"] = plan.Description.ValueString()
	}
	if !plan.WindowType.IsUnknown() {
		evaluation["window_type"] = plan.WindowType.ValueString()
	}
	if !plan.WindowDurationMin.IsUnknown() {
		evaluation["window_duration_minutes"] = plan.WindowDurationMin.ValueInt64()
	}
	if !plan.GroupBy.IsNull() {
		evaluation["group_by"] = StringListValueToStringSlice(plan.GroupBy)
	}
	if !plan.EventTimestamp.IsNull() {
		evaluation["event_timestamp"] = plan.EventTimestamp.ValueString()
	}
	if !plan.Script.IsNull() {
		evaluation["script"] = plan.Script.ValueString()
	}

	return &component, dd
}

// From an API response to a terraform model
func ChangeAlertToModel(plan *ChangeAlertModel, component *Alert) {
	plan.Id = NewStringValue(component.Id)
	plan.ComponentKind = NewStringValue(component.ComponentKind)
	plan.ComponentId = NewStringValue(component.ComponentId)
	plan.Active = NewBoolValue(component.Active)
	if component.Inputs != nil {
		inputs := make([]attr.Value, 0)
		for _, v := range component.Inputs {
			inputs = append(inputs, NewStringValue(v))
		}
		plan.Inputs = NewListValueMust(StringType{}, inputs)
	}

	general := component.AlertConfig["general"].(map[string]any)
	evaluation := component.AlertConfig["evaluation"].(map[string]any)

	// General properties
	plan.Name = NewStringValue(general["name"].(string))
	if general["description"] != nil {
		plan.Description = NewStringValue(general["description"].(string))
	}

	// Evaluation properties
	plan.EventType = NewStringValue(evaluation["event_type"].(string))
	if evaluation["group_by"] != nil {
		plan.GroupBy = SliceToStringListValue(evaluation["group_by"].([]any))
	}
	plan.Operation = NewStringValue(FindKey(Aggregate_Operations, evaluation["operation"].(string)))
	plan.WindowType = NewStringValue(evaluation["window_type"].(string))

	plan.WindowDurationMin = NewInt64Value(int64(evaluation["window_duration_minutes"].(float64)))
	plan.Conditional = UnwindConditionalToModel(evaluation["conditional"].(map[string]any), Non_Change_Operator_Labels)
	if evaluation["script"] != nil {
		plan.Script = NewStringValue(evaluation["script"].(string))
	}
	if evaluation["event_timestamp"] != nil {
		plan.EventTimestamp = NewStringValue(evaluation["event_timestamp"].(string))
	}

	// Alert Payload properties
	plan.AlertPayload = GetAlertPayloadToModel(component.AlertConfig["alert_payload"].(map[string]any))
}
