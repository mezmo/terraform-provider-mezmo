package processors

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	. "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
)

const REDUCE_PROCESSOR_NODE_NAME = "reduce"
const REDUCE_PROCESSOR_TYPE_NAME = REDUCE_PROCESSOR_NODE_NAME

type ReduceProcessorModel struct {
	Id              StringValue `tfsdk:"id"`
	PipelineId      StringValue `tfsdk:"pipeline_id"`
	Title           StringValue `tfsdk:"title"`
	Description     StringValue `tfsdk:"description"`
	Inputs          ListValue   `tfsdk:"inputs"`
	GenerationId    Int64Value  `tfsdk:"generation_id"`
	DurationMs      Int64Value  `tfsdk:"duration_ms" user_config:"true"`
	MaxEvents       Int64Value  `tfsdk:"max_events" user_config:"true"`
	GroupBy         ListValue   `tfsdk:"group_by" user_config:"true"`
	DateFormats     ListValue   `tfsdk:"date_formats" user_config:"true"`
	MergeStrategies ListValue   `tfsdk:"merge_strategies" user_config:"true"`
	FlushCondition  ObjectValue `tfsdk:"flush_condition" user_config:"true"`
}

var ReduceProcessorResourceSchema = schema.Schema{
	Description: "Combine multiple events over time into one based on a set of criteria",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"duration_ms": schema.Int64Attribute{
			Optional: true,
			Description: "The amount of time (in milliseconds) to allow streaming events to accumulate " +
				"into a single \"reduced\" event. The process repeats indefinitely, or until " +
				"an \"ends when\" condition is satisfied.",
			Computed: true,
			Default:  int64default.StaticInt64(30000),
		},
		"max_events": schema.Int64Attribute{
			Optional: true,
			Description: "The maximum number of events that can be included in a time window (specified " +
				"by duration_ms). The reduce operation will stop once it has reached this " +
				"number of events, regardless of whether the duration_ms have elapsed.",
		},
		"group_by": schema.ListAttribute{
			ElementType: StringType{},
			Optional:    true,
			Description: "Before reducing, group events based on matching data from each of these " +
				"field paths. Supports nesting via dot-notation.",
		},
		"date_formats": schema.ListNestedAttribute{
			Optional: true,
			Description: "Describes which root-level properties are dates, and their expected format. " +
				"Dot-notation is supported, but nested field lookup paths will be an error.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"field": schema.StringAttribute{
						Required:    true,
						Description: "Specifies a root-level path property that contains a date value.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							stringvalidator.LengthAtMost(200),
						},
					},
					"format": schema.StringAttribute{
						Required:    true,
						Description: "The template describing the date format",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							stringvalidator.LengthAtMost(200),
						},
					},
				},
			},
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
			},
		},
		"merge_strategies": schema.ListNestedAttribute{
			Optional: true,
			Description: "Specify merge strategies for individual root-level properties. " +
				"Dot-notation is supported, but nested field lookup paths will be an error.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"field": schema.StringAttribute{
						Required:    true,
						Description: "This is a root-level path property to apply a merge strategy to its value",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							stringvalidator.LengthAtMost(200),
						},
					},
					"strategy": schema.StringAttribute{
						Required:    true,
						Description: "The merge strategy to be used for the specified property",
						Validators: []validator.String{
							stringvalidator.OneOf(ReduceMergeStrategies...),
						},
					},
				},
			},
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
			},
		},
		"flush_condition": schema.SingleNestedAttribute{
			Optional: true,
			Description: "Force accumulated event reduction to flush the result when a " +
				"conditional expression evaluates to true on an inbound event.",
			Attributes: map[string]schema.Attribute{
				"when": schema.StringAttribute{
					Required: true,
					Description: "Specifies whether to start a new reduction of events based on the " +
						"conditions, or end a current reduction based on them.",
					Validators: []validator.String{
						stringvalidator.OneOf("starts_when", "ends_when"),
					},
				},
				"conditional": ParentConditionalAttribute,
			},
		},
	}),
}

func ReduceProcessorFromModel(plan *ReduceProcessorModel, previousState *ReduceProcessorModel) (*Processor, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Processor{
		BaseNode: BaseNode{
			Type:        REDUCE_PROCESSOR_NODE_NAME,
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig: map[string]any{
				"duration_ms": plan.DurationMs.ValueInt64(),
			},
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	if !plan.Inputs.IsUnknown() {
		inputs := make([]string, 0)
		for _, v := range plan.Inputs.Elements() {
			value, _ := v.(StringValue)
			inputs = append(inputs, value.ValueString())
		}
		component.Inputs = inputs
	}

	if !plan.MaxEvents.IsNull() {
		component.UserConfig["max_events"] = plan.MaxEvents.ValueInt64()
	}
	if !plan.GroupBy.IsNull() {
		component.UserConfig["group_by"] = StringListValueToStringSlice(plan.GroupBy)
	}

	if !plan.DateFormats.IsNull() {
		dateFormats := make([]map[string]any, 0)
		for _, v := range plan.DateFormats.Elements() {
			obj := MapValuesToMapAny(v, &dd)
			dateFormats = append(dateFormats, obj)
		}
		component.UserConfig["date_formats"] = dateFormats
	}

	if !plan.MergeStrategies.IsNull() {
		mergeStrategies := make([]map[string]any, 0)
		for _, v := range plan.MergeStrategies.Elements() {
			obj := MapValuesToMapAny(v, &dd)
			mergeStrategies = append(mergeStrategies, obj)
		}
		component.UserConfig["merge_strategies"] = mergeStrategies
	}

	if !plan.FlushCondition.IsNull() {
		flushCondition := plan.FlushCondition.Attributes()
		component.UserConfig["flush_condition"] = map[string]any{
			"when":        GetAttributeValue[StringValue](flushCondition, "when").ValueString(),
			"conditional": unwindConditionalFromModel(flushCondition["conditional"]),
		}
	}

	return &component, dd
}

func unwindConditionalFromModel(v attr.Value) map[string]any {
	conditional := map[string]any{
		"expressions":       nil,
		"logical_operation": "AND",
	}
	value, ok := v.(ObjectValue)
	if !ok {
		panic(fmt.Errorf("Expected an object but did not receive one: %+v", v))
	}
	attrs := value.Attributes()

	if !attrs["logical_operation"].IsUnknown() {
		conditional["logical_operation"] = attrs["logical_operation"].(StringValue).ValueString()
	}

	if expressions, ok := attrs["expressions"]; ok && !expressions.IsNull() {
		// Loop and extract expressions into an array of map[strings]
		elements := expressions.(ListValue).Elements()
		result := make([]map[string]any, 0, len(elements))
		for _, obj := range elements {
			propVals := obj.(ObjectValue).Attributes()
			elem := map[string]any{
				"field":        propVals["field"].(StringValue).ValueString(),
				"str_operator": propVals["operator"].(StringValue).ValueString(),
			}
			if propVals["value_number"].IsNull() {
				elem["value"] = propVals["value_string"].(StringValue).ValueString()
			} else {
				elem["value"] = propVals["value_number"].(Float64Value).ValueFloat64()
			}
			result = append(result, elem)
		}
		conditional["expressions"] = result

	} else if group, ok := attrs["expressions_group"]; ok && !group.IsNull() {
		// This is a nested list of conditionals. Recursively unwind them.
		conditionals := group.(ListValue).Elements()
		result := make([]map[string]any, 0, len(conditionals))
		for _, conditional := range conditionals {
			result = append(result, unwindConditionalFromModel(conditional))
		}
		conditional["expressions"] = result
	}

	return conditional
}

func ReduceProcessorToModel(plan *ReduceProcessorModel, component *Processor) {
	// plan.ClearFields()
	plan.Id = NewStringValue(component.Id)
	if component.Title != "" {
		plan.Title = NewStringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = NewStringValue(component.Description)
	}
	plan.GenerationId = NewInt64Value(component.GenerationId)
	if component.Inputs != nil {
		inputs := make([]attr.Value, 0)
		for _, v := range component.Inputs {
			inputs = append(inputs, NewStringValue(v))
		}
		plan.Inputs = NewListValueMust(StringType{}, inputs)
	}

	plan.DurationMs = NewInt64Value(int64(component.UserConfig["duration_ms"].(float64)))

	if component.UserConfig["max_events"] != nil {
		plan.MaxEvents = NewInt64Value(int64(component.UserConfig["max_events"].(float64)))
	}
	if component.UserConfig["group_by"] != nil {
		plan.GroupBy = SliceToStringListValue(component.UserConfig["group_by"].([]any))
	}

	if component.UserConfig["date_formats"] != nil {
		dateFormats := make([]attr.Value, 0)
		for _, v := range component.UserConfig["date_formats"].([]any) {
			dateFormats = append(dateFormats, NewObjectValueMust(map[string]attr.Type{
				"field":  StringType{},
				"format": StringType{},
			}, map[string]attr.Value{
				"field":  NewStringValue(v.(map[string]any)["field"].(string)),
				"format": NewStringValue(v.(map[string]any)["format"].(string)),
			}))
		}
		plan.DateFormats = NewListValueMust(ObjectType{
			AttrTypes: dateFormats[0].Type(context.Background()).(ObjectType).AttributeTypes(),
		}, dateFormats)
	}

	if component.UserConfig["merge_strategies"] != nil {
		mergeStrategies := make([]attr.Value, 0)
		for _, v := range component.UserConfig["merge_strategies"].([]any) {
			mergeStrategies = append(mergeStrategies, NewObjectValueMust(map[string]attr.Type{
				"field":    StringType{},
				"strategy": StringType{},
			}, map[string]attr.Value{
				"field":    NewStringValue(v.(map[string]any)["field"].(string)),
				"strategy": NewStringValue(v.(map[string]any)["strategy"].(string)),
			}))
		}
		plan.MergeStrategies = NewListValueMust(ObjectType{
			AttrTypes: mergeStrategies[0].Type(context.Background()).(ObjectType).AttributeTypes(),
		}, mergeStrategies)
	}

	if component.UserConfig["flush_condition"] != nil {
		flushCondition := component.UserConfig["flush_condition"].(map[string]any)
		whenValue := flushCondition["when"].(string)
		if whenValue == "starts_when" || whenValue == "ends_when" {
			conditional := UnwindConditionalToModel(flushCondition["conditional"].(map[string]any))
			plan.FlushCondition = NewObjectValueMust(map[string]attr.Type{
				"when":        StringType{},
				"conditional": conditional.Type(context.Background()),
			}, map[string]attr.Value{
				"when":        NewStringValue(whenValue),
				"conditional": conditional,
			})
		}
	}
}
