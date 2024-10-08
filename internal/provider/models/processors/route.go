package processors

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/models/modelutils"
)

type RouteProcessorModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline_id"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	Inputs       List   `tfsdk:"inputs"`
	GenerationId Int64  `tfsdk:"generation_id"`
	Conditionals List   `tfsdk:"conditionals" user_config:"true"`
	// unmatched exists in outputs from API which is not exposed to the
	// user. Mapped to user_config to make TF happy
	Unmatched String `tfsdk:"unmatched" user_config:"true"`
}

const ROUTE_PROCESSOR_TYPE_NAME = "route"
const ROUTE_PROCESSOR_NODE_NAME = ROUTE_PROCESSOR_TYPE_NAME

var RouteProcessorResourceSchema = schema.Schema{
	Description: "Route data based on whether or not it matches logical comparisons.",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"conditionals": schema.ListNestedAttribute{
			Required:    true,
			Description: "A list of conditions, each of which has a label and an expression or expression groups.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: ExtendSchemaAttributes(ParentConditionalAttribute(Non_Change_Operator_Labels).Attributes, map[string]schema.Attribute{
					"label": schema.StringAttribute{
						Required:    true,
						Description: "A label for the expression group",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							stringvalidator.LengthAtMost(255),
						},
					},
					"output_name": schema.StringAttribute{
						Computed: true,
						Description: "A system generated value to identify the results of this expression. " +
							"This value should be used when connecting the results to another processor or destination.",
					},
				}),
			},
		},
		"unmatched": schema.StringAttribute{
			Computed:    true,
			Description: "A system generated value to identify the results that don't match any condition.",
		},
	}),
}

func RouteProcessorFromModel(plan *RouteProcessorModel, previousState *RouteProcessorModel) (*Processor, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Processor{
		BaseNode: BaseNode{
			Type:        ROUTE_PROCESSOR_TYPE_NAME,
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	component.Inputs = StringListValueToStringSlice(plan.Inputs)

	component.UserConfig = conditionalsFromModel(plan.Conditionals)
	return &component, dd
}

func RouteProcessorToModel(plan *RouteProcessorModel, component *Processor) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}
	plan.GenerationId = Int64Value(component.GenerationId)
	plan.Inputs = SliceToStringListValue(component.Inputs)

	outputs := mappedOutputs(component)
	if unmatched, ok := outputs["_unmatched"]; ok {
		plan.Unmatched = unmatched
	}

	conditionals, ok := component.UserConfig["conditionals"].([]any)
	if ok {
		elemType := plan.Conditionals.ElementType(context.Background())
		if elemType == nil {
			// used by ConvertToTerraformModel method
			// when hydrating without a terraform state, the list element type is nil
			listType := RouteProcessorResourceSchema.Attributes["conditionals"].GetType()
			elemType = listType.(basetypes.ListType).ElementType()
		}

		list_value, diag := conditionalsToModel(conditionals, elemType, outputs)
		if !diag.HasError() {
			plan.Conditionals = list_value
		}
	}
}

func conditionalsFromModel(v List) map[string]any {
	var conditionals []map[string]any
	for _, entry := range v.Elements() {
		conditionals = append(conditionals, map[string]any{
			"conditional": UnwindConditionalFromModel(entry),
			"label":       entry.(Object).Attributes()["label"].(String).ValueString(),
		})
	}

	return map[string]any{
		"conditionals": conditionals,
	}
}

func conditionalsToModel(respConditionals []any, listItemType attr.Type, outputs map[string]basetypes.StringValue) (List, diag.Diagnostics) {
	var conditionals []basetypes.ObjectValue

	for _, entry := range respConditionals {
		conditional := entry.(map[string]any)["conditional"].(map[string]any)
		unwound := UnwindConditionalToModel(conditional, Non_Change_Operator_Labels)

		attrTypes := unwound.AttributeTypes(context.Background())
		attrTypes["label"] = basetypes.StringType{}
		attrTypes["output_name"] = basetypes.StringType{}

		attrValues := unwound.Attributes()
		attrValues["label"] = StringValue(entry.(map[string]any)["label"].(string))
		apiOutputName := entry.(map[string]any)["_output_name"].(string)
		if apiOutputName != "" {
			if outputName, ok := outputs[strings.ToLower(apiOutputName)]; ok {
				attrValues["output_name"] = outputName
			}
		}

		unwound = basetypes.NewObjectValueMust(attrTypes, attrValues)
		conditionals = append(conditionals, unwound)
	}

	return basetypes.NewListValueFrom(context.Background(), listItemType, conditionals)
}
