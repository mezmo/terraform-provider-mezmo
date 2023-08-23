package modelutils

import (
	"reflect"

	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// This function can help with taking an object from an API response and translating
// it into a basetype object from the terraform-plugin-framework
func MapStringsToMapValues(values map[string]string) map[string]attr.Value {
	result := make(map[string]attr.Value, len(values))
	for k, v := range values {
		result[k] = StringValue(v)
	}
	return result
}

// This function can receive a terraform object (as defined in their basetypes) and
// return a regular map[string]string object that can be used in `Component` such that
// it can be used in an API call.
// TODO: This currently only supports string values, not numbers or any other type.
// TODO: perhaps this should accept a mutable interface that can be inspected for casting
func MapValuesToMapStrings(obj interface{}, dd diag.Diagnostics) (map[string]string, bool) {
	var attrs map[string]attr.Value
	switch obj.(type) {
	case Object:
		attrs = obj.(Object).Attributes()
	case Map:
		attrs = obj.(Map).Elements()
	default:
		panic(fmt.Sprintf("Unsupported object type for `fromAttributes`: %s", reflect.TypeOf(obj)))
	}
	target := make(map[string]string, len(attrs))
	for k, v := range attrs {
		if v.IsUnknown() {
			continue
		}
		stringValue, ok := attrs[k].(basetypes.StringValue)
		if !ok {
			dd.AddError(
				"Could not look up attribute value",
				fmt.Sprintf("Cannot cast key %s to a string value. Please report this to Mezmo.", k),
			)
			continue
		}
		target[k] = stringValue.ValueString()
	}
	return target, dd.HasError()
}

// Gets a known attribute value from an attribute map and casts it to the provided type.
func GetAttributeValue[T any](m map[string]attr.Value, name string) T {
	r, _ := m[name].(T)
	return r
}

func StringListValueToStringSlice(list List) []string {
	if list.IsUnknown() {
		return nil
	}

	result := make([]string, 0)
	for _, v := range list.Elements() {
		value, _ := v.(basetypes.StringValue)
		result = append(result, value.ValueString())
	}

	return result
}

func SliceToStringListValue[T any](s []T) List {
	list := make([]attr.Value, 0)
	for _, v := range s {
		value, _ := any(v).(string)
		list = append(list, StringValue(value))
	}
	return ListValueMust(StringType, list)
}
