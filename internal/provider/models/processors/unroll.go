package processors

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/models/modelutils"
)

type UnrollProcessorModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline_id"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	Inputs       List   `tfsdk:"inputs"`
	GenerationId Int64  `tfsdk:"generation_id"`
	Field        String `tfsdk:"field" user_config:"true"`
	ValuesOnly   Bool   `tfsdk:"values_only" user_config:"true"`
}

func UnrollProcessorResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Takes an array of events and emits them all as individual events",
		Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
			"field": schema.StringAttribute{
				Required:    true,
				Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
				Description: "The field name that contains an array of events",
			},
			"values_only": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
				Description: "When enabled, the values from the specified array field will be emitted as " +
					"new events. Otherwise, the original event will be duplicated for each value " +
					"in the array field, with the unrolled value present in the field specified.",
			},
		}),
	}
}

func UnrollProcessorFromModel(plan *UnrollProcessorModel, previousState *UnrollProcessorModel) (*Processor, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Processor{
		BaseNode: BaseNode{
			Type:        "unroll",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			Inputs:      StringListValueToStringSlice(plan.Inputs),
			UserConfig: map[string]any{
				"field":       plan.Field.ValueString(),
				"values_only": plan.ValuesOnly.ValueBool(),
			},
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component, dd
}

func UnrollProcessorToModel(plan *UnrollProcessorModel, component *Processor) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}
	plan.GenerationId = Int64Value(component.GenerationId)
	plan.Inputs = SliceToStringListValue(component.Inputs)
	plan.Field = StringValue(component.UserConfig["field"].(string))
	plan.ValuesOnly = BoolValue(component.UserConfig["values_only"].(bool))
}
