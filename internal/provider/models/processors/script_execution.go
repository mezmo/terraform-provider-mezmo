package processors

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
)

type ScriptExecutionProcessorModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline_id"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	Inputs       List   `tfsdk:"inputs"`
	GenerationId Int64  `tfsdk:"generation_id"`
	Script       String `tfsdk:"script" user_config:"true"`
}

var ScriptExecutionProcessorResourceSchema = schema.Schema{
	Description: "Use JavaScript to reshape and transform your data" +
		" You can combine multiple actions like filtering, dropping," +
		" mapping, and casting inside of a single js script",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"script": schema.StringAttribute{
			Required:   true,
			Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			Description: "The script containing the JavaScript function that represents the " +
				"transformation of events flowing though the pipeline",
		},
	}),
}

func ScriptExecutionProcessorFromModel(plan *ScriptExecutionProcessorModel, previousState *ScriptExecutionProcessorModel) (*Processor, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Processor{
		BaseNode: BaseNode{
			Type:        "js-script",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			Inputs:      StringListValueToStringSlice(plan.Inputs),
			UserConfig: map[string]any{
				"script": plan.Script.ValueString(),
			},
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component, dd
}

func ScriptExecutionProcessorToModel(plan *ScriptExecutionProcessorModel, component *Processor) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}
	plan.GenerationId = Int64Value(component.GenerationId)
	plan.Inputs = SliceToStringListValue(component.Inputs)
	plan.Script = StringValue(component.UserConfig["script"].(string))
}
