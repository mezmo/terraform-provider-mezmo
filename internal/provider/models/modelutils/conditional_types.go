package modelutils

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var ExpressionTypes = map[string]attr.Type{
	"field":        StringType{},
	"operator":     StringType{},
	"value_number": Float64Type{},
	"value_string": StringType{},
}

var ExpressionAttributes = map[string]schema.Attribute{
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

var LogicalOperationAttribute = schema.StringAttribute{
	Optional:    true,
	Computed:    true,
	Description: "The logical operation (AND/OR) to be applied to the list of conditionals",
	Validators: []validator.String{
		stringvalidator.OneOf("AND", "OR"),
	},
}

var ExpressionListAttribute = schema.ListNestedAttribute{
	Optional:    true,
	Description: "Defines a list of expressions for field comparisons",
	NestedObject: schema.NestedAttributeObject{
		Attributes: ExpressionAttributes,
	},
	Validators: []validator.List{
		listvalidator.SizeAtLeast(1),
		listvalidator.ConflictsWith(
			path.MatchRelative().AtParent().AtName("expressions_group"),
		),
	},
}

var NestedExpressionAttribute = schema.NestedAttributeObject{
	Attributes: ExpressionAttributes,
}

var MAX_NESTED_LEVELS = 5
var childExpressionGroups = make([]schema.NestedAttributeObject, MAX_NESTED_LEVELS+1, MAX_NESTED_LEVELS+1)

func getChildExpressionGroupAttributes(level int) schema.NestedAttributeObject {
	// TODO: The addition of `route` can easily add a `label` attribute here
	attributes := map[string]schema.Attribute{
		"expressions":       ExpressionListAttribute,
		"logical_operation": LogicalOperationAttribute,
	}

	if level < MAX_NESTED_LEVELS {
		attributes["expressions_group"] = schema.ListNestedAttribute{
			Optional:     true,
			Description:  "A group of expressions joined by a logical operator",
			NestedObject: getChildExpressionGroupAttributes(level + 1),
		}
	}

	nestedAttribute := schema.NestedAttributeObject{
		Attributes: attributes,
	}
	childExpressionGroups[level] = nestedAttribute

	return nestedAttribute
}

var ParentConditionalAttribute = schema.SingleNestedAttribute{
	Optional:    true,
	Description: "A group of expressions (optionally nested) joined by a logical operator",
	Attributes: map[string]schema.Attribute{
		"expressions": ExpressionListAttribute,
		"expressions_group": schema.ListNestedAttribute{
			Optional:     true,
			Description:  "A group of expressions joined by a logical operator",
			NestedObject: getChildExpressionGroupAttributes(0),
		},
		"logical_operation": LogicalOperationAttribute,
	},
}

func GetAttributesByLevel(level int) map[string]schema.Attribute {
	if level == 0 {
		return ParentConditionalAttribute.Attributes
	}
	return childExpressionGroups[level-1].Attributes // Zero-based array, so subtract 1
}

func GetChildExpressionGroupTypeByLevel(level int) attr.Type {
	return childExpressionGroups[level].Type()
}
