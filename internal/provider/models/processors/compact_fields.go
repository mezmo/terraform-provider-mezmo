package processors

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
	"github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
)

type CompactFieldsProcessorModel struct {
	Id            String `tfsdk:"id"`
	PipelineId    String `tfsdk:"pipeline_id"`
	Title         String `tfsdk:"title"`
	Description   String `tfsdk:"description"`
	Inputs        List   `tfsdk:"inputs"`
	GenerationId  Int64  `tfsdk:"generation_id"`
	Fields        List   `tfsdk:"fields" user_config:"true"`
	CompactArray  Bool   `tfsdk:"compact_array" user_config:"true"`
	CompactObject Bool   `tfsdk:"compact_object" user_config:"true"`
}

var CompactFieldsProcessorResourceSchema = schema.Schema{
	Description: "Remove empty values from a list of fields",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"fields": schema.ListAttribute{
			ElementType: StringType,
			Required:    true,
			Description: "A list of fields to remove empty values from",
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
				listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
			},
		},
		"compact_array": schema.BoolAttribute{
			Optional:    true,
			Description: "Remove empty arrays from a field",
			Computed:    true,
			Default:     booldefault.StaticBool(true),
		},
		"compact_object": schema.BoolAttribute{
			Optional:    true,
			Description: "Remove empty objects from a field",
			Computed:    true,
			Default:     booldefault.StaticBool(true),
		},
	}),
}

func CompactFieldsProcessorFromModel(plan *CompactFieldsProcessorModel, previousState *CompactFieldsProcessorModel) (*Processor, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Processor{
		BaseNode: BaseNode{
			Type:        "compact-fields",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig:  make(map[string]any),
		},
	}

	var options = make(map[string]bool)
	component.UserConfig["options"] = options

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
		options["compact_array"] = previousState.CompactArray.ValueBool()
		options["compact_object"] = previousState.CompactObject.ValueBool()
	}

	if !plan.Inputs.IsUnknown() {
		inputs := make([]string, 0)
		for _, v := range plan.Inputs.Elements() {
			value, _ := v.(basetypes.StringValue)
			inputs = append(inputs, value.ValueString())
		}
		component.Inputs = inputs
	}

	component.UserConfig["fields"] = modelutils.StringListValueToStringSlice(plan.Fields)

	if !plan.CompactArray.IsUnknown() {
		options["compact_array"] = plan.CompactArray.ValueBool()
	}
	if !plan.CompactObject.IsUnknown() {
		options["compact_object"] = plan.CompactObject.ValueBool()
	}

	return &component, dd
}

func CompactFieldsProcessorToModel(plan *CompactFieldsProcessorModel, component *Processor) {
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

	plan.Fields = modelutils.SliceToStringListValue(component.UserConfig["fields"].([]any))
	options, _ := component.UserConfig["options"].(map[string]any)
	if options != nil {
		plan.CompactArray = BoolValue(options["compact_array"].(bool))
		plan.CompactObject = BoolValue(options["compact_object"].(bool))
	}
}
