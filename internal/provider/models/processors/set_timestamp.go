package processors

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/v5/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/v5/internal/provider/models/modelutils"
)

const SET_TIMESTAMP_PROCESSOR_TYPE_NAME = "set_timestamp"
const SET_TIMESTAMP_PROCESSOR_NODE_NAME = "set-timestamp"

type SetTimestampProcessorModel struct {
	Id           StringValue `tfsdk:"id"`
	PipelineId   StringValue `tfsdk:"pipeline_id"`
	Title        StringValue `tfsdk:"title"`
	Description  StringValue `tfsdk:"description"`
	Inputs       ListValue   `tfsdk:"inputs"`
	GenerationId Int64Value  `tfsdk:"generation_id"`
	Parsers      ListValue   `tfsdk:"parsers" user_config:"true"`
}

var parserAttrTypes = map[string]attr.Type{
	"field":            StringType{},
	"timestamp_format": StringType{},
}

var SetTimestampProcessorResourceSchema = schema.Schema{
	Description: "Parse a list of fields using a chosen time format. First match will override the default timestamp of the event.",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"parsers": schema.ListNestedAttribute{
			Required:    true,
			Description: "The list of fields and parsers to use in order of priority, short-circuiting on the first successful match.",
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
			},
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"field": schema.StringAttribute{
						Required:    true,
						Description: "The field whose value will be parsed.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"timestamp_format": schema.StringAttribute{
						Required:    true,
						Description: "The stftime timestamp format.  This will be stored as a custom format in this processor.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
			},
		},
	}),
}

func SetTimestampProcessorFromModel(plan *SetTimestampProcessorModel, previousState *SetTimestampProcessorModel) (*Processor, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Processor{
		BaseNode: BaseNode{
			Type:        SET_TIMESTAMP_PROCESSOR_NODE_NAME,
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig:  map[string]any{},
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}
	component.Inputs = StringListValueToStringSlice(plan.Inputs)

	var parsers []map[string]any
	for _, v := range plan.Parsers.Elements() {
		parserItem := make(map[string]any)
		optionsItem := make(map[string]string)

		parserPlanMap := MapValuesToMapAny(v, &dd)
		parserItem["field"] = parserPlanMap["field"]
		optionsItem["format"] = "Custom"
		optionsItem["custom_format"] = parserPlanMap["timestamp_format"].(string)
		parserItem["options"] = optionsItem
		parsers = append(parsers, parserItem)
	}
	component.UserConfig["parsers"] = parsers

	return &component, dd
}

func SetTimestampProcessorToModel(plan *SetTimestampProcessorModel, component *Processor) {
	plan.Id = NewStringValue(component.Id)
	if component.Title != "" {
		plan.Title = NewStringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = NewStringValue(component.Description)
	}
	plan.GenerationId = NewInt64Value(component.GenerationId)
	plan.Inputs = SliceToStringListValue(component.Inputs)

	parsers := make([]attr.Value, 0)
	for _, v := range component.UserConfig["parsers"].([]any) {
		parser := v.(map[string]any)

		optionsArray, _ := parser["options"].(map[string]any)

		attrValues := map[string]attr.Value{
			"field":            types.StringNull(),
			"timestamp_format": types.StringNull(),
		}

		attrValues["field"] = types.StringValue(parser["field"].(string))
		attrValues["timestamp_format"] = types.StringValue(optionsArray["custom_format"].(string))

		parsers = append(parsers, NewObjectValueMust(parserAttrTypes, attrValues))
	}
	plan.Parsers = NewListValueMust(ObjectType{AttrTypes: parserAttrTypes}, parsers)

}
