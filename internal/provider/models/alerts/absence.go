package alerts

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	. "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
)

const ALERT_TYPE_ABSENCE = "absence"

type AbsenceAlertModel struct {
	Id                StringValue `tfsdk:"id"`
	PipelineId        StringValue `tfsdk:"pipeline_id"`
	ComponentKind     StringValue `tfsdk:"component_kind"`
	ComponentId       StringValue `tfsdk:"component_id"`
	Inputs            ListValue   `tfsdk:"inputs"`
	Active            BoolValue   `tfsdk:"active"`
	Name              StringValue `tfsdk:"name" user_config:"true"`
	Description       StringValue `tfsdk:"description" user_config:"true"`
	EventType         StringValue `tfsdk:"event_type" user_config:"true"`
	GroupBy           ListValue   `tfsdk:"group_by" user_config:"true"`
	WindowType        StringValue `tfsdk:"window_type" user_config:"true"`
	WindowDurationMin Int64Value  `tfsdk:"window_duration_minutes" user_config:"true"`
	EventTimestamp    StringValue `tfsdk:"event_timestamp" user_config:"true"`
	AlertPayload      ObjectValue `tfsdk:"alert_payload" user_config:"true"`
}

var AbsenceAlertResourceSchema = schema.Schema{
	Description: "Represents an Absence Alert in a Pipeline",
	Attributes:  baseAlertSchemaAttributes,
}

// From terraform schema/model to a struct for sending to the API
func AbsenceAlertFromModel(plan *AbsenceAlertModel, previousState *AbsenceAlertModel) (*Alert, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	alertPayload := GetAlertPayloadFromModel(plan.AlertPayload, &dd)

	CustomErrorChecks(&CheckedFields{
		EventType:      plan.EventType,
		EventTimestamp: plan.EventTimestamp,
		GroupBy:        plan.GroupBy,
	}, &dd)

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
				"alert_type": ALERT_TYPE_ABSENCE, // Required for the API, but hidden from the user here
				"event_type": plan.EventType.ValueString(),
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

	return &component, dd
}

// From an API response to a terraform model
func AbsenceAlertToModel(plan *AbsenceAlertModel, component *Alert) {
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
	plan.WindowType = NewStringValue(evaluation["window_type"].(string))

	plan.WindowDurationMin = NewInt64Value(int64(evaluation["window_duration_minutes"].(float64)))
	if evaluation["event_timestamp"] != nil {
		plan.EventTimestamp = NewStringValue(evaluation["event_timestamp"].(string))
	}

	// Alert Payload properties
	plan.AlertPayload = GetAlertPayloadToModel(component.AlertConfig["alert_payload"].(map[string]any))
}
