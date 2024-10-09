package modelutils

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var ExpressionTypes = map[string]attr.Type{
	"field":        StringType{},
	"operator":     StringType{},
	"value_number": Float64Type{},
	"value_string": StringType{},
	"negate":       BoolType{},
}

func expressionAttributes(operators []string) map[string]schema.Attribute {
	possibleValues := strings.Join(operators[:len(operators)-1], ", ") + " or " + operators[len(operators)-1]

	return map[string]schema.Attribute{
		"field": schema.StringAttribute{
			Required:    true,
			Description: "The field path whose value will be used in the comparison",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"operator": schema.StringAttribute{
			Required:    true,
			Description: "The comparison operator. Possible values are: " + possibleValues + ".",
			Validators: []validator.String{
				stringvalidator.OneOf(operators...),
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
		"negate": schema.BoolAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Negate the operator",
			Validators: []validator.Bool{
				boolvalidator.AlsoRequires(
					path.MatchRelative().AtParent().AtName("operator"),
				),
			},
		},
	}
}

var LogicalOperationAttribute = schema.StringAttribute{
	Optional:    true,
	Computed:    true,
	Description: "The logical operation (AND/OR) to be applied to the list of conditionals",
	Validators: []validator.String{
		stringvalidator.OneOf("AND", "OR"),
	},
}

func expressionListAttribute(operators []string) schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional:    true,
		Description: "Defines a list of expressions for field comparisons",
		NestedObject: schema.NestedAttributeObject{
			Attributes: expressionAttributes(operators),
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
			listvalidator.ConflictsWith(
				path.MatchRelative().AtParent().AtName("expressions_group"),
			),
		},
	}
}

func nestedExpressionAttribute(operators []string) schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: expressionAttributes(operators),
	}
}

var MAX_NESTED_LEVELS = 5
var childExpressionGroups = make([]schema.NestedAttributeObject, MAX_NESTED_LEVELS+1, MAX_NESTED_LEVELS+1)

func getChildExpressionGroupAttributes(level int, operators []string) schema.NestedAttributeObject {
	// TODO: The addition of `route` can easily add a `label` attribute here
	attributes := map[string]schema.Attribute{
		"expressions":       expressionListAttribute(operators),
		"logical_operation": LogicalOperationAttribute,
	}

	if level < MAX_NESTED_LEVELS {
		attributes["expressions_group"] = schema.ListNestedAttribute{
			Optional:     true,
			Description:  "A group of expressions joined by a logical operator",
			NestedObject: getChildExpressionGroupAttributes(level+1, operators),
		}
	}

	nestedAttribute := schema.NestedAttributeObject{
		Attributes: attributes,
	}
	childExpressionGroups[level] = nestedAttribute

	return nestedAttribute
}

func ParentConditionalAttribute(operators []string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:    true,
		Description: "A group of expressions (optionally nested) joined by a logical operator",
		Attributes: map[string]schema.Attribute{
			"expressions": expressionListAttribute(operators),
			"expressions_group": schema.ListNestedAttribute{
				Optional:     true,
				Description:  "A group of expressions joined by a logical operator",
				NestedObject: getChildExpressionGroupAttributes(0, operators),
			},
			"logical_operation": LogicalOperationAttribute,
		},
	}
}

func getAttributesByLevel(level int, operators []string) map[string]schema.Attribute {
	if level == 0 {
		return ParentConditionalAttribute(operators).Attributes
	}
	return childExpressionGroups[level-1].Attributes // Zero-based array, so subtract 1
}

func getChildExpressionGroupTypeByLevel(level int) attr.Type {
	return childExpressionGroups[level].Type()
}

func UnwindConditionalToModel(component map[string]any, operators []string) ObjectValue {
	value, _ := parseExpressionsItem(component, 0, operators)
	return value
}

func parseExpressionsItem(component map[string]any, level int, operators []string) (value ObjectValue, isGroup bool) {
	logicalOperation := "AND" // Default
	if operation, ok := component["logical_operation"].(string); ok {
		logicalOperation = operation
	}
	if childExpressionArr, ok := component["expressions"].([]any); ok {
		// Branch
		groupItems := make([]attr.Value, 0)
		leafItems := make([]attr.Value, 0)
		attributeTypes := ToAttrTypes(getAttributesByLevel(level, operators))

		for _, e := range childExpressionArr {
			child := e.(map[string]any)
			value, isGroup := parseExpressionsItem(child, level+1, operators)

			if isGroup {
				groupItems = append(groupItems, value)
			} else {
				leafItems = append(leafItems, value)
			}
		}

		attributeValues := map[string]attr.Value{
			"logical_operation": NewStringValue(logicalOperation),
			"expressions":       NewListNull(nestedExpressionAttribute(operators).Type()), // Default to match `plan` since there might only be group expressions on this level
		}

		expressionGroupType := getChildExpressionGroupTypeByLevel(level)
		mergeLeavesIntoGroups := len(leafItems) > 0 && len(groupItems) > 0
		if mergeLeavesIntoGroups {
			aTypes := expressionGroupType.(basetypes.ObjectType).AttributeTypes()
			expressions := NewObjectValueMust(aTypes, map[string]attr.Value{
				"logical_operation": NewStringValue("AND"),
				"expressions":       NewListValueMust(nestedExpressionAttribute(operators).Type(), leafItems),
				"expressions_group": NewListNull(aTypes["expressions_group"].(basetypes.ListType).ElementType()),
			})
			groupItems = append([]attr.Value{expressions}, groupItems...)
		} else if len(leafItems) > 0 {
			attributeValues["expressions"] = NewListValueMust(nestedExpressionAttribute(operators).Type(), leafItems)
		}

		if len(groupItems) > 0 {
			attributeValues["expressions_group"] = NewListValueMust(expressionGroupType, groupItems)
		} else if attributeTypes["expressions_group"] != nil {
			attributeValues["expressions_group"] = NewListNull(expressionGroupType)
		}

		return NewObjectValueMust(attributeTypes, attributeValues), true
	}

	// Leaf
	attributeValues := map[string]attr.Value{
		"field":        NewStringValue(component["field"].(string)),
		"operator":     NewStringValue(component["str_operator"].(string)),
		"value_number": NewFloat64Null(),
		"value_string": NewStringNull(),
		"negate":       NewBoolValue(false),
	}
	if negate, ok := component["negate"].(bool); ok {
		attributeValues["negate"] = NewBoolValue(negate)
	}
	if valueNumber, ok := component["value"].(float64); ok {
		attributeValues["value_number"] = NewFloat64Value(valueNumber)
	} else if valueString, ok := component["value"].(string); ok && valueString != "" {
		attributeValues["value_string"] = NewStringValue(valueString)
	}

	return NewObjectValueMust(ExpressionTypes, attributeValues), false
}

func UnwindConditionalFromModel(v attr.Value) map[string]any {
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
			if propVals["negate"].IsNull() {
				elem["negate"] = NewBoolValue(false)
			} else {
				elem["negate"] = propVals["negate"].(BoolValue).ValueBool()
			}
			result = append(result, elem)
		}
		conditional["expressions"] = result

	} else if group, ok := attrs["expressions_group"]; ok && !group.IsNull() {
		// This is a nested list of conditionals. Recursively unwind them.
		conditionals := group.(ListValue).Elements()
		result := make([]map[string]any, 0, len(conditionals))
		for _, conditional := range conditionals {
			result = append(result, UnwindConditionalFromModel(conditional))
		}
		conditional["expressions"] = result
	}

	return conditional
}
