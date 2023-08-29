package transforms

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
	"github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/models/modelutils"
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
			stringvalidator.OneOf(modelutils.Operators...),
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
}

var logicalOperation = schema.StringAttribute{
	Optional:    true,
	Description: "The logical operation (AND/OR) to be applied to the list of conditionals",
	Validators: []validator.String{
		stringvalidator.OneOf("AND", "OR"),
	},
}

var expressionList = schema.ListNestedAttribute{
	Required:    true,
	Description: "Defines an expression for field comparison",
	NestedObject: schema.NestedAttributeObject{
		Attributes: expressionAttributes,
	},
	Validators: []validator.List{
		listvalidator.SizeAtLeast(1),
	},
}

var expressionGroup = schema.ListNestedAttribute{
	Optional:    true,
	Description: "A group of expressions joined by a logical operator",
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"expressions":       expressionList,
			"logical_operation": logicalOperation, // Needs function for optional
		},
	},
	// function should take validators for conflicts with `expressions`
}

var conditionalValueType = schema.SingleNestedAttribute{
	Optional:    true,
	Description: "xxx",
	Attributes: map[string]schema.Attribute{
		"expressions":       expressionList,  // Needs function for required optional
		"expressions_group": expressionGroup, // *** Needs function to add another expression_group
		"logical_operation": logicalOperation,
	},
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
								stringvalidator.OneOf(modelutils.ReduceMergeStrategies...),
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
						Required:    true,
						Description: "The merge strategy to be used for the specified property",
						Validators: []validator.String{
							stringvalidator.OneOf("starts_when", "ends_when"),
						},
					},
				},
			},
		}),
	}
}

func ReduceTransformFromModel(plan *ReduceTransformModel, previousState *ReduceTransformModel) (*Transform, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Transform{
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
	}

	if !plan.Inputs.IsUnknown() {
		inputs := make([]string, 0)
		for _, v := range plan.Inputs.Elements() {
			value, _ := v.(basetypes.StringValue)
			inputs = append(inputs, value.ValueString())
		}
		component.Inputs = inputs
	}

	return &component, dd
}

func ReduceTransformToModel(plan *ReduceTransformModel, component *Transform) {
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
}
