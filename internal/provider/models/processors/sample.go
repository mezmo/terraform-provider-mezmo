package processors

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
)

type SampleProcessorModel struct {
	Id            String `tfsdk:"id"`
	PipelineId    String `tfsdk:"pipeline_id"`
	Title         String `tfsdk:"title"`
	Description   String `tfsdk:"description"`
	Inputs        List   `tfsdk:"inputs"`
	GenerationId  Int64  `tfsdk:"generation_id"`
	Rate          Int64  `tfsdk:"rate" user_config:"true"`
	AlwaysInclude Object `tfsdk:"always_include" user_config:"true"`
}

var SampleProcessorResourceSchema = schema.Schema{
	Description: "Sample data at a given rate, retaining only a subset of data events for further processing",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"rate": schema.Int64Attribute{
			Computed: true,
			Optional: true,
			Description: "The rate at which events will be forwarded, expressed as 1/N. For example," +
				" `rate = 10` means 1 out of every 10 events will be forwarded and the rest" +
				" will be dropped",
			Validators: []validator.Int64{
				int64validator.AtLeast(2),
				int64validator.AtMost(10000),
			},
			Default: int64default.StaticInt64(10),
		},
		"always_include": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Events matching this criteria will always show up in the results",
			Attributes: map[string]schema.Attribute{
				"field": schema.StringAttribute{
					Required:    true,
					Description: "The field to use in a condition to always include in sampling",
					Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
				},
				"operator": schema.StringAttribute{
					Required: true,
					Description: "The comparison operator to check the value of the field or" +
						" whether the first exists",
					Validators: []validator.String{
						stringvalidator.OneOf(Operators...),
					},
				},
				"value_string": schema.StringAttribute{
					Optional:    true,
					Description: "The operand to compare the field value with, when the value is a string",
					Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
				},
				"value_number": schema.Float64Attribute{
					Optional:    true,
					Description: "The operand to compare the field value with, when the value is a number",
				},
				"case_sensitive": schema.BoolAttribute{
					Optional: true,
					Computed: true,
					// Default:     booldefault.StaticBool(true),
					Description: "Perform case sensitive comparison?",
					Validators: []validator.Bool{
						boolvalidator.AlsoRequires(
							path.MatchRelative().AtParent().AtName("operator"),
							path.MatchRelative().AtParent().AtName("value_string"),
						),
					},
				},
			},
		},
	}),
}

func SampleProcessorFromModel(plan *SampleProcessorModel, previousState *SampleProcessorModel) (*Processor, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Processor{
		BaseNode: BaseNode{
			Type:        "sample",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig: map[string]any{
				"rate": plan.Rate.ValueInt64(),
			},
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	component.Inputs = StringListValueToStringSlice(plan.Inputs)

	if !plan.AlwaysInclude.IsNull() {
		plan_map := plan.AlwaysInclude.Attributes()
		component_map := make(map[string]any)
		component_map["field"] = GetAttributeValue[String](plan_map, "field").ValueString()
		component_map["str_operator"] = GetAttributeValue[String](plan_map, "operator").ValueString()
		if !plan_map["value_string"].IsNull() {
			component_map["value"] = GetAttributeValue[String](plan_map, "value_string").ValueString()
		} else if !plan_map["value_number"].IsNull() {
			component_map["value"] = GetAttributeValue[Float64](plan_map, "value_number").ValueFloat64()
		}

		op := component_map["str_operator"]
		if op == "equal" || op == "no_equal" || op == "starts_with" || op == "ends_with" {
			component_map["case_sensitive"] = GetAttributeValue[Bool](plan_map, "case_sensitive").ValueBool()
		}

		component.UserConfig["always_include"] = component_map
	}

	return &component, dd
}

func SampleProcessorToModel(plan *SampleProcessorModel, component *Processor) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}
	plan.GenerationId = Int64Value(component.GenerationId)
	plan.Inputs = SliceToStringListValue(component.Inputs)
	plan.Rate = Int64Value(int64(component.UserConfig["rate"].(float64)))
	if component.UserConfig["always_include"] != nil {
		component_map, _ := component.UserConfig["always_include"].(map[string]any)
		if len(component_map) > 0 {
			plan_map := make(map[string]attr.Value)
			plan_map["field"] = StringValue(component_map["field"].(string))
			plan_map["operator"] = StringValue(component_map["str_operator"].(string))
			plan_map["value_number"] = Float64Null()
			plan_map["value_string"] = StringNull()
			plan_map["case_sensitive"] = BoolNull()
			if value, ok := component_map["value"]; ok {
				if valueString, ok := value.(string); ok {
					if valueString != "" {
						plan_map["value_string"] = StringValue(valueString)
					}
				} else if valueNumber, ok := value.(float64); ok {
					plan_map["value_number"] = Float64Value(valueNumber)
				} else {
					panic("Unexpected type for value in always_include field")
				}
			}
			if case_sensitive, ok := component_map["case_sensitive"]; ok {
				plan_map["case_sensitive"] = BoolValue(case_sensitive.(bool))
			}
			objT := plan.AlwaysInclude.AttributeTypes(context.Background())
			if len(objT) == 0 {
				objT = SampleProcessorResourceSchema.Attributes["always_include"].GetType().(basetypes.ObjectType).AttrTypes
			}
			plan.AlwaysInclude = basetypes.NewObjectValueMust(objT, plan_map)
		}
	}
}
