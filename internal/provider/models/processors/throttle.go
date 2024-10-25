package processors

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/models/modelutils"
)

const THROTTLE_PROCESSOR_NODE_NAME = "throttle"
const THROTTLE_PROCESSOR_TYPE_NAME = THROTTLE_PROCESSOR_NODE_NAME

type ThrottleProcessorModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline_id"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	Inputs       List   `tfsdk:"inputs"`
	GenerationId Int64  `tfsdk:"generation_id"`
	Threshold    Int64  `tfsdk:"threshold" user_config:"true"`
	WindowMS     Int64  `tfsdk:"window_ms" user_config:"true"`
	KeyField     String `tfsdk:"key_field" user_config:"true"`
	Exclude      Object `tfsdk:"exclude" user_config:"true"`
}

var ThrottleProcessorResourceSchema = schema.Schema{
	Description: "Throttle (rate-limit) events passing through this component",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"threshold": schema.Int64Attribute{
			Required:    true,
			Description: "The number of events to allow over the time window. If a Key Field is specified, the limit will be applied to events for each unique value of that field",
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},
		"window_ms": schema.Int64Attribute{
			Required:    true,
			Description: "The time window over which the configured limit is applied",
			Validators: []validator.Int64{
				int64validator.AtLeast(1000),
			},
		},
		"key_field": schema.StringAttribute{
			Optional:    true,
			Description: "The field to use as a key for rate limiting. If specified, the rate limit will be applied to each unique value of this field",
		},
		"exclude": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Events matching this criteria will be excluded from the rate limit",
			Attributes:  ParentConditionalAttribute(Non_Change_Operator_Labels).Attributes,
		},
	}),
}

func ThrottleProcessorFromModel(plan *ThrottleProcessorModel, previousState *ThrottleProcessorModel) (*Processor, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Processor{
		BaseNode: BaseNode{
			Type:        THROTTLE_PROCESSOR_NODE_NAME,
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig: map[string]any{
				"threshold": plan.Threshold.ValueInt64(),
				"window_ms": plan.WindowMS.ValueInt64(),
			},
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	component.Inputs = StringListValueToStringSlice(plan.Inputs)

	if !plan.KeyField.IsNull() {
		component.UserConfig["key_field"] = plan.KeyField.ValueString()
	}

	if !plan.Exclude.IsNull() {
		component.UserConfig["exclude"] = UnwindConditionalFromModel(plan.Exclude)
	}

	return &component, dd
}

func ThrottleProcessorToModel(plan *ThrottleProcessorModel, component *Processor) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}

	plan.GenerationId = Int64Value(component.GenerationId)
	plan.Inputs = SliceToStringListValue(component.Inputs)
	plan.Threshold = Int64Value(int64(component.UserConfig["threshold"].(float64)))
	plan.WindowMS = Int64Value(int64(component.UserConfig["window_ms"].(float64)))

	if keyField, ok := component.UserConfig["key_field"]; ok {
		plan.KeyField = StringValue(keyField.(string))
	}

	if exclude, ok := component.UserConfig["exclude"]; ok {
		plan.Exclude = UnwindConditionalToModel(exclude.(map[string]any), Non_Change_Operator_Labels)
	}
}
