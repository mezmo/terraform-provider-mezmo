package processors

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/v5/internal/client"
)

const STRINGIFY_PROCESSOR_NODE_NAME = "stringify"
const STRINGIFY_PROCESSOR_TYPE_NAME = STRINGIFY_PROCESSOR_NODE_NAME

type StringifyProcessorModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline_id"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	Inputs       List   `tfsdk:"inputs"`
	GenerationId Int64  `tfsdk:"generation_id"`
}

var StringifyProcessorResourceSchema = schema.Schema{
	Description: "Represents a processor to stringify JSON data.",
	Attributes:  ExtendBaseAttributes(map[string]schema.Attribute{}),
}

func StringifyProcessorFromModel(plan *StringifyProcessorModel, previousState *StringifyProcessorModel) (*Processor, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Processor{
		BaseNode: BaseNode{
			Type:        STRINGIFY_PROCESSOR_NODE_NAME,
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig:  make(map[string]any),
		},
	}

	if !plan.Inputs.IsUnknown() {
		inputs := make([]string, 0)
		for _, v := range plan.Inputs.Elements() {
			value, _ := v.(basetypes.StringValue)
			inputs = append(inputs, value.ValueString())
		}
		component.Inputs = inputs
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component, dd
}

func StringifyProcessorToModel(plan *StringifyProcessorModel, component *Processor) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}
	plan.GenerationId = Int64Value(component.GenerationId)
	if component.Inputs != nil {
		inputs := make([]attr.Value, 0)
		for _, v := range component.Inputs {
			inputs = append(inputs, StringValue(v))
		}
		plan.Inputs = ListValueMust(StringType, inputs)
	}
}
