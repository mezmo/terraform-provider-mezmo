package processors

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	types "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/mezmo/terraform-provider-mezmo/v4/internal/client"
	modelutils "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/models/modelutils"
)

const TRACE_SAMPLING_PROCESSOR_NODE_NAME = "trace-sampling"
const TRACE_SAMPLING_PROCESSOR_TYPE_NAME = "trace_sampling"

type TraceSamplingProcessorModel struct {
	Id                types.String `tfsdk:"id"`
	PipelineId        types.String `tfsdk:"pipeline_id"`
	Title             types.String `tfsdk:"title"`
	Description       types.String `tfsdk:"description"`
	Inputs            types.List   `tfsdk:"inputs"`
	GenerationId      types.Int64  `tfsdk:"generation_id"`
	SampleType        types.String `tfsdk:"sample_type" user_config:"true"`
	TraceIdField      types.String `tfsdk:"trace_id_field" user_config:"true"`
	Rate              types.Int64  `tfsdk:"rate" user_config:"true"`
	ParentSpanIdField types.String `tfsdk:"parent_span_id_field" user_config:"true"`
	Conditionals      types.List   `tfsdk:"conditionals" user_config:"true"`
}

var TraceSamplingProcessorResourceSchema = schema.Schema{
	Description: "The Trace Sampling Processor allows you to sample traces using either 'head' or 'tail' based methodologies.",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"sample_type": schema.StringAttribute{
			Required:    true,
			Description: "The type of sampling to apply. Can be one of the following: 'head' or 'tail'.",
			Validators: []validator.String{
				stringvalidator.OneOf("head", "tail"),
			},
		},
		"trace_id_field": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The field name of the trace ID to sample on.",
		},
		"rate": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Head Sampling: The rate at which to sample traces expressed as 1/N. A value of 2 means 50% sampling, 3 means 33%, etc.",
			Validators: []validator.Int64{
				int64validator.AtLeast(2),
				int64validator.AtMost(100000),
			},
		},
		"parent_span_id_field": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Tail Sampling: The field name to pull the parent span id. Absence of this field indicates the event is a root span.",
		},
		"conditionals": schema.ListNestedAttribute{
			Optional:    true,
			Description: "Tail Sampling: A list of complete conditionals, including expressions and a sample rate to maintain.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: (map[string]schema.Attribute{
					"conditional": schema.SingleNestedAttribute{
						Computed:    true,
						Optional:    true,
						Description: "Tail Sampling: A conditional expression for sampling.",
						Attributes:  ExtendSchemaAttributes(modelutils.ParentConditionalAttribute(modelutils.Non_Change_Operator_Labels).Attributes, map[string]schema.Attribute{}),
					},
					"rate": schema.Int64Attribute{
						Required:    true,
						Description: "The rate at which to sample traces expressed as 1/N. A value of 1 means 100% sampling, 2 means 50%, 3 means 33%, etc.",
						Validators: []validator.Int64{
							int64validator.AtLeast(1),
							int64validator.AtMost(100000),
						},
					},
					"_output_name": schema.StringAttribute{
						Computed:    true,
						Optional:    true,
						Description: "The generated name of the rule.",
					},
				}),
			},
		},
	}),
}

func TraceSamplingProcessorFromModel(plan *TraceSamplingProcessorModel, previousState *TraceSamplingProcessorModel) (*client.Processor, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	// a lot of config props are marked optional, so we need to check if they are set
	// because they may not actually be optional or applicable given the type
	if !ValidateRequiredConfigsPerSampleType(plan, plan.SampleType.ValueString(), &dd) {
		return nil, dd
	}

	component := client.Processor{
		BaseNode: client.BaseNode{
			Type:        TRACE_SAMPLING_PROCESSOR_NODE_NAME,
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
		},
	}
	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}
	component.Inputs = modelutils.StringListValueToStringSlice(plan.Inputs)
	component.UserConfig = map[string]any{
		"sample_type": plan.SampleType.ValueString(),
	}
	if !plan.TraceIdField.IsNull() && len(plan.TraceIdField.ValueString()) > 0 {
		component.UserConfig["trace_id_field"] = plan.TraceIdField.ValueString()
	}

	if plan.SampleType.ValueString() == "head" {
		if !plan.Rate.IsNull() && plan.Rate.ValueInt64() != 0 {
			component.UserConfig["rate"] = plan.Rate.ValueInt64()
		}
	} else {
		if !plan.ParentSpanIdField.IsNull() && len(plan.ParentSpanIdField.ValueString()) > 0 {
			component.UserConfig["parent_span_id_field"] = plan.ParentSpanIdField.ValueString()
		}

		var conditionals []map[string]any
		for _, entry := range plan.Conditionals.Elements() {
			conditionals = append(conditionals, map[string]any{
				"conditional":  modelutils.UnwindConditionalFromModel(entry.(basetypes.ObjectValue).Attributes()["conditional"]),
				"rate":         entry.(basetypes.ObjectValue).Attributes()["rate"].(basetypes.Int64Value).ValueInt64(),
				"_output_name": entry.(basetypes.ObjectValue).Attributes()["_output_name"].(basetypes.StringValue).ValueString(),
			})
		}

		component.UserConfig["conditionals"] = conditionals
	}

	return &component, dd
}

func TraceSamplingProcessorToModel(plan *TraceSamplingProcessorModel, component *client.Processor) {
	plan.Id = types.StringValue(component.Id)
	if component.Title != "" {
		plan.Title = types.StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = types.StringValue(component.Description)
	}
	plan.GenerationId = types.Int64Value(component.GenerationId)
	plan.Inputs = modelutils.SliceToStringListValue(component.Inputs)
	plan.SampleType = types.StringValue(component.UserConfig["sample_type"].(string))
	plan.TraceIdField = types.StringValue(component.UserConfig["trace_id_field"].(string))

	if plan.SampleType.ValueString() == "head" {
		if rate, ok := component.UserConfig["rate"].(float64); ok {
			plan.Rate = types.Int64Value(int64(rate))
		}

		plan.ParentSpanIdField = types.StringNull()
	} else if plan.SampleType.ValueString() == "tail" {
		plan.Rate = types.Int64Null()

		if parent_span_id_field, ok := component.UserConfig["parent_span_id_field"].(string); ok {
			plan.ParentSpanIdField = types.StringValue(parent_span_id_field)
		}

		if conditionals, ok := component.UserConfig["conditionals"].([]any); ok {
			elemType := plan.Conditionals.ElementType(context.Background())
			if elemType == nil {
				// used by ConvertToTerraformModel method
				// when hydrating without a terraform state, the list element type is nil
				listType := TraceSamplingProcessorResourceSchema.Attributes["conditionals"].GetType()
				elemType = listType.(basetypes.ListType).ElementType()
			}

			list_value, diag := traceSamplingConditionalsToModel(conditionals, elemType)
			if !diag.HasError() {
				plan.Conditionals = list_value
			}
		}
	}
}

func traceSamplingConditionalsToModel(respConditionals []any, listItemType attr.Type) (types.List, diag.Diagnostics) {
	var conditionals []basetypes.ObjectValue

	for _, entry := range respConditionals {
		var conditional = entry.(map[string]any)["conditional"].(map[string]any)
		var rate float64 = entry.(map[string]any)["rate"].(float64)
		var outputName string = entry.(map[string]any)["_output_name"].(string)

		var unwound = modelutils.UnwindConditionalToModel(conditional, modelutils.Non_Change_Operator_Labels)

		var attrTypes = map[string]attr.Type{}
		attrTypes["rate"] = basetypes.Int64Type{}
		attrTypes["conditional"] = unwound.Type(context.Background())
		attrTypes["_output_name"] = basetypes.StringType{}

		var attrValues = map[string]attr.Value{}
		attrValues["rate"] = types.Int64Value(int64(rate))
		attrValues["conditional"] = unwound
		attrValues["_output_name"] = types.StringValue(outputName)

		var result = basetypes.NewObjectValueMust(attrTypes, attrValues)
		conditionals = append(conditionals, result)
	}

	return basetypes.NewListValueFrom(context.Background(), listItemType, conditionals)
}

func ValidateRequiredConfigsPerSampleType(plan *TraceSamplingProcessorModel, sample_type string, dd *diag.Diagnostics) bool {
	result := true
	if sample_type == "head" {
		if !plan.ParentSpanIdField.IsNull() && len(plan.ParentSpanIdField.ValueString()) > 0 {
			dd.AddAttributeError(
				path.Root("parent_span_id_field"),
				"Attribute \"parent_span_id_field\" is not applicable for head sampling.",
				"Attribute \"parent_span_id_field\" is not applicable for head sampling.",
			)
			result = false
		}
		if !plan.Conditionals.IsNull() {
			dd.AddAttributeError(
				path.Root("conditionals"),
				"Attribute \"conditionals\" is not applicable for head sampling.",
				"Attribute \"conditionals\" is not applicable for head sampling.",
			)
			result = false
		}
	} else if sample_type == "tail" {
		if !plan.Rate.IsNull() && plan.Rate.ValueInt64() != 0 {
			dd.AddAttributeError(
				path.Root("rate"),
				"Attribute \"rate\" is not applicable for tail sampling.",
				"Attribute \"rate\" is not applicable for tail sampling.",
			)
			result = false
		}
		if plan.Conditionals.IsNull() || len(plan.Conditionals.Elements()) == 0 {
			dd.AddAttributeError(
				path.Root("conditionals"),
				"Attribute \"conditionals\" is required for tail sampling.",
				"Attribute \"conditionals\" is required for tail sampling.",
			)
			result = false
		}
	}

	return result
}
