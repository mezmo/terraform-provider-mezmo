package transforms

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
)

type StringifyTransformModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	Inputs       List   `tfsdk:"inputs"`
	GenerationId Int64  `tfsdk:"generation_id"`
}

func StringifyTransformResourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"pipeline": schema.StringAttribute{
				Required:    true,
				Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
				Description: "The pipeline identifier",
			},
			"title": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(256),
				},
			},
			"inputs": schema.ListAttribute{
				ElementType: StringType,
				Optional:    true,
				Description: "The ids of the input components",
			},
			"generation_id": schema.Int64Attribute{
				Computed: true,
			},
		},
	}
}

func StringifyTransformFromModel(model *StringifyTransformModel, previousState *StringifyTransformModel) *Component {
	component := Component{
		Type:        "demo-logs",
		Title:       model.Title.ValueString(),
		Description: model.Description.ValueString(),
		UserConfig:  make(map[string]any),
	}

	if !model.Inputs.IsUnknown() {
		inputs := make([]string, 0)
		for _, v := range model.Inputs.Elements() {
			value, _ := v.(basetypes.StringValue)
			inputs = append(inputs, value.ValueString())
		}
		component.Inputs = inputs
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component
}

func StringifyTransformToModel(model *StringifyTransformModel, component *Component) {
	model.Id = StringValue(component.Id)
	if component.Title != "" {
		model.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		model.Description = StringValue(component.Description)
	}
	model.GenerationId = Int64Value(component.GenerationId)
	if component.Inputs != nil {
		inputs := make([]attr.Value, 0)
		for _, v := range component.Inputs {
			inputs = append(inputs, StringValue(v))
		}
		model.Inputs = ListValueMust(StringType, inputs)
	}
}
