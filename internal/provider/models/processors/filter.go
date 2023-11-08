package processors

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	. "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
)

const FILTER_PROCESSOR_NODE_NAME = "filter"
const FILTER_PROCESSOR_TYPE_NAME = "filter"

type FilterProcessorModel struct {
	Id           StringValue `tfsdk:"id"`
	PipelineId   StringValue `tfsdk:"pipeline_id"`
	Title        StringValue `tfsdk:"title"`
	Description  StringValue `tfsdk:"description"`
	Inputs       ListValue   `tfsdk:"inputs"`
	GenerationId Int64Value  `tfsdk:"generation_id"`
	Action       StringValue `tfsdk:"action" user_config:"true"`
	Conditional  ObjectValue `tfsdk:"conditional" user_config:"true"`
}

var FilterProcessorResourceSchema = schema.Schema{
	Description: "Define condition(s) to include or exclude events from the pipeline",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"action": schema.StringAttribute{
			Description: "How to handle events matching this criteria",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.OneOf("allow", "drop"),
			},
		},
		"conditional": schema.SingleNestedAttribute{
			Required:    true,
			Description: ParentConditionalAttribute.Description,
			Attributes:  ParentConditionalAttribute.Attributes,
		},
	}),
}

func FilterProcessorFromModel(plan *FilterProcessorModel, previousState *FilterProcessorModel) (*Processor, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := &Processor{
		BaseNode: BaseNode{
			Type:        FILTER_PROCESSOR_NODE_NAME,
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig: map[string]any{
				"action": plan.Action.ValueString(),
			},
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	component.Inputs = StringListValueToStringSlice(plan.Inputs)
	component.UserConfig["conditional"] = unwindConditionalFromModel(plan.Conditional)

	return component, dd
}

func FilterProcessorToModel(plan *FilterProcessorModel, component *Processor) {
	plan.Id = NewStringValue(component.Id)
	if component.Title != "" {
		plan.Title = NewStringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = NewStringValue(component.Description)
	}
	plan.GenerationId = NewInt64Value(component.GenerationId)
	plan.Inputs = SliceToStringListValue(component.Inputs)
	plan.Action = NewStringValue(component.UserConfig["action"].(string))
	plan.Conditional = unwindConditionalToModel(component.UserConfig["conditional"].(map[string]any))
}
