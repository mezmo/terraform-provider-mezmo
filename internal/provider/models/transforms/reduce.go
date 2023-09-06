package transforms

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/models/modelutils"
)

type ReduceTransformModel struct {
	Id              String `tfsdk:"id"`
	PipelineId      String `tfsdk:"pipeline_id"`
	Title           String `tfsdk:"title"`
	Description     String `tfsdk:"description"`
	Inputs          List   `tfsdk:"inputs"`
	GenerationId    Int64  `tfsdk:"generation_id"`
	DurationMs      Int64  `tfsdk:"duration_ms"`
	GroupBy         List   `tfsdk:"group_by"`
	DateFormats     List   `tfsdk:"date_formats"`
	MergeStrategies List   `tfsdk:"merge_strategies"`
	FlushCondition  Object `tfsdk:"flush_condition"`
}

var expressionTypes = map[string]attr.Type{
	"field":        StringType,
	"operator":     StringType,
	"value_number": Float64Type,
	"value_string": StringType,
}

var expressionAttributes = map[string]schema.Attribute{
	"field": schema.StringAttribute{
		Required:    true,
		Description: "The field path whose value will be used in the comparison",
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	},
	"operator": schema.StringAttribute{
		Required:    true,
		Description: "The comparison operator",
		Validators: []validator.String{
			stringvalidator.OneOf(Operators...),
		},
	},
	"value_string": schema.StringAttribute{
		Optional:    true,
		Description: "The operand to compare the field value with, when the value is a string",
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
			stringvalidator.ConflictsWith(
				path.MatchRelative().AtParent().AtName("value_number"),
			),
		},
	},
	"value_number": schema.Float64Attribute{
		Optional:    true,
		Description: "The operand to compare the field value with, when the value is a number",
	},
}

var logicalOperation = schema.StringAttribute{
	Optional:    true,
	Computed:    true,
	Description: "The logical operation (AND/OR) to be applied to the list of conditionals",
	Validators: []validator.String{
		stringvalidator.OneOf("AND", "OR"),
	},
}

var expressionList = schema.ListNestedAttribute{
	Optional:    true,
	Description: "Defines a list of expressions for field comparisons",
	NestedObject: expressionItem,
	Validators: []validator.List{
		listvalidator.SizeAtLeast(1),

		// WOW: TIL about ConflictsWith()
		listvalidator.ConflictsWith(
			path.MatchRelative().AtParent().AtName("expressions_group"),
		),
	},
}

var expressionItem = schema.NestedAttributeObject{
	Attributes: expressionAttributes,
}

var conditionalValueType = schema.SingleNestedAttribute{
	Optional:    true,
	Description: "A group of expressions (optionally nested) joined by a logical operator",
	Attributes: map[string]schema.Attribute{
		"expressions":       expressionList,
		"expressions_group": schema.ListNestedAttribute{
			Optional:    true,
			Description: "A group of nested expressions joined by a logical operator",
			NestedObject: nestedObjectL1,
		},
		"logical_operation": logicalOperation,
	},
}

var nestedObjectL1 = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"expressions":       expressionList,
		"expressions_group": schema.ListNestedAttribute{
			Optional:    true,
			Description: "A group of nested expressions joined by a logical operator",
			NestedObject: nestedObjectL2,
		},
		"logical_operation": logicalOperation,
	},
}

var nestedObjectL2 = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"expressions":       expressionList,
		"expressions_group": schema.ListNestedAttribute{
			Optional:    true,
			Description: "A group of nested expressions joined by a logical operator",
			NestedObject: nestedObjectL3,
		},
		"logical_operation": logicalOperation,
	},
}

var nestedObjectL3 = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"expressions":       expressionList,
		"logical_operation": logicalOperation,
	},
}

var attributesByLevel = []map[string]schema.Attribute{
	conditionalValueType.Attributes, nestedObjectL1.Attributes, nestedObjectL2.Attributes,
}

var childExpressionGroupTypeByLevel = []attr.Type{
	nestedObjectL1.Type(), nestedObjectL2.Type(), nestedObjectL3.Type(),
}

func ReduceTransformResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Remove empty values from a list of fields",
		Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
			"duration_ms": schema.Int64Attribute{
				Optional: true,
				Description: "The amount of time (in milliseconds) to allow streaming events to accumulate " +
					"into a single \"reduced\" event. The process repeats indefinitely, or until " +
					"an \"ends when\" condition is satisfied.",
				Computed: true,
				Default:  int64default.StaticInt64(30000),
			},
			"group_by": schema.ListAttribute{
				ElementType: StringType,
				Optional:    true,
				Description: "Before reducing, group events based on matching data from each of these " +
					"field paths. Supports nesting via dot-notation.",
			},
			"date_formats": schema.ListNestedAttribute{
				Optional: true,
				Description: "Before reducing, group events based on matching data from each of these " +
					"field paths. Supports nesting via dot-notation.",
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
					"conditional": conditionalValueType,
				},
			},
		}),
	}
}

func ReduceTransformFromModel(plan *ReduceTransformModel, previousState *ReduceTransformModel) (*Transform, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Transform{
		BaseNode: BaseNode{
			Type:        "reduce",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig: map[string]any{
				"duration_ms": plan.DurationMs.ValueInt64(),
			},
		},
	}

	fmt.Printf("---------------- %+v\n", plan)

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	if !plan.Inputs.IsUnknown() {
		inputs := make([]string, 0)
		for _, v := range plan.Inputs.Elements() {
			value, _ := v.(basetypes.StringValue)
			inputs = append(inputs, value.ValueString())
		}
		component.Inputs = inputs
	}

	if !plan.FlushCondition.IsNull() {
		flushCondition := plan.FlushCondition.Attributes()
		component.UserConfig["flush_condition"] = map[string]any{
			"when":        GetAttributeValue[String](flushCondition, "when").ValueString(),
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
	value, ok := v.(basetypes.ObjectValue)
	if !ok {
		panic(fmt.Errorf("Expected an object but did not receive one: %+v", v))
	}
	attrs := value.Attributes()

	if !attrs["logical_operation"].IsNull() {
		conditional["logical_operation"] = attrs["logical_operation"].(String).ValueString()
	}

	if expressions, ok := attrs["expressions"]; ok && !expressions.IsNull() {
		// Loop and extract expressions into an array of map[strings]
		elements := expressions.(List).Elements()
		result := make([]map[string]any, 0, len(elements))
		for _, obj := range elements {
			propVals := obj.(basetypes.ObjectValue).Attributes()
			elem := map[string]any{
				"field":        propVals["field"].(String).ValueString(),
				"str_operator": propVals["operator"].(String).ValueString(),
			}
			if propVals["value_number"].IsNull() {
				elem["value"] = propVals["value_string"].(String).ValueString()
			} else {
				elem["value"] = propVals["value_number"].(Float64).ValueFloat64()
			}
			result = append(result, elem)
		}
		conditional["expressions"] = result

	} else if group, ok := attrs["expressions_group"]; ok && !group.IsNull() {
		// This is a nested list of conditionals. Recursively unwind them.
		conditionals := group.(basetypes.ListValue).Elements()
		result := make([]map[string]any, 0, len(conditionals))
		for _, conditional := range conditionals {
			result = append(result, unwindConditionalFromModel(conditional))
		}
		conditional["expressions"] = result
	}

	return conditional
}

func ReduceTransformToModel(plan *ReduceTransformModel, component *Transform) {
	PrintJSON(component)
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
	if component.UserConfig["flush_condition"] != nil {
		flushCondition := component.UserConfig["flush_condition"].(map[string]any)
		// TODO: set conditional
		conditional := unwindConditionalToModel(flushCondition["conditional"].(map[string]any))
		fmt.Println("--conditional", conditional)
	}
}

func unwindConditionalToModel(component map[string]any) basetypes.ObjectValue {
	incoming, _ := json.MarshalIndent(component, "", "   ")
	fmt.Println("--unwinding", string(incoming))

	value, _ := parseExpressionsItem(component, component["logical_operation"].(string), 0)

	return value
}

func parseExpressionsItem(component map[string]any, logicalOperation string, level int) (value basetypes.ObjectValue, isGroup bool) {
	if childExpressionArr, ok := component["expressions"].([]any); ok {
		// Branch
		groupItems := make([]attr.Value, 0)
		leafItems := make([]attr.Value, 0)
		attributeTypes := toAttrTypes(attributesByLevel[level])

		for _, e := range childExpressionArr {
			value, isGroup := parseExpressionsItem(
				e.(map[string]any),
				component["logical_operation"].(string),
				level+1)

			if isGroup {
				groupItems = append(groupItems, value)
			} else {
				leafItems = append(leafItems, value)
			}
		}

		attributeValues := map[string]attr.Value{
			"logical_operation": StringValue(logicalOperation),
			"expressions": ListValueMust(expressionItem.Type(), leafItems),
		}

		expressionGroupType := childExpressionGroupTypeByLevel[level]
		if len(groupItems) > 0 {
			attributeValues["expressions_group"] = ListValueMust(expressionGroupType, groupItems)
		} else if attributeTypes["expressions_group"] != nil {
			attributeValues["expressions_group"] = ListNull(expressionGroupType)
		}

		return basetypes.NewObjectValueMust(
			attributeTypes,
			attributeValues), true
	}

	// Leaf
	fmt.Println("----LEAF", level, component)
	// It's an attributeValues object to covert to a TF type
	attributeValues := map[string]attr.Value{
		"field":        StringValue(component["field"].(string)),
		"operator":     StringValue(component["str_operator"].(string)),
		"value_number": Float64Null(),
		"value_string": StringNull(),
	}
	if valueNumber, ok := component["value"].(float64); ok {
		attributeValues["value_number"] = basetypes.NewFloat64Value(valueNumber)
	} else {
		attributeValues["value_string"] = basetypes.NewStringValue(component["value"].(string))
	}

	return basetypes.NewObjectValueMust(expressionTypes, attributeValues), false
}

func toAttrTypes(attributes map[string]schema.Attribute) map[string]attr.Type{
	result := make(map[string]attr.Type)
	for k, v := range attributes {
		result[k] = v.GetType()
	}

	return result
}
