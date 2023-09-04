package modelutils

import (
	"reflect"

	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"golang.org/x/exp/slices"
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

// Convert from an API response object to a terraform plugin framework
// basetype object
func MapAnyToMapValues(attrs map[string]attr.Type, values map[string]any, optional_fields []string) map[string]attr.Value {
	result := make(map[string]attr.Value, len(values))
	for k, v := range values {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Bool:
			result[k] = BoolValue(v.(bool))
		case reflect.Float32, reflect.Float64:
			result[k] = Float64Value(v.(float64))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			result[k] = Int64Value(v.(int64))
		case reflect.String:
			if slices.Contains(optional_fields, k) && (v == nil || v.(string) == "") {
				result[k] = StringNull()
			} else {
				result[k] = StringValue(v.(string))
			}
		case reflect.Slice:
			if slices.Contains(optional_fields, k) && v == nil {
				result[k] = ListNull(StringType)
			} else {
				result[k] = SliceToStringListValue(v.([]string))
			}
		}
	}

	PopulateMissingMapValues(attrs, result)
	return result
}

func PopulateMissingMapValues(attrs map[string]attr.Type, result map[string]attr.Value) {
	for field, attr_type := range attrs {
		if _, ok := result[field]; ok {
			continue
		}
		switch attr_type {
		case basetypes.BoolType{}:
			result[field] = BoolNull()
		case basetypes.Int64Type{}:
			result[field] = Int64Null()
		case basetypes.Float64Type{}:
			result[field] = Float64Null()
		case basetypes.NumberType{}:
			result[field] = NumberNull()
		case basetypes.StringType{}:
			result[field] = StringNull()
		case basetypes.ListType{ElemType: basetypes.StringType{}}:
			result[field] = ListNull(basetypes.StringType{})
		}
	}
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

// Returns the keys of the provided map
func MapKeys[K comparable, V any](m map[K]V) []K {
	result := make([]K, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

// Returns the key for a given map, panicking when not found
func FindKey[K comparable, V comparable](m map[K]V, value V) K {
	for k, v := range m {
		if v == value {
			return k
		}
	}
	panic(fmt.Sprintf("Key for '%v' not found", value))
}

// Sets one or more optional string values from a attribute map into a target map
func SetOptionalStringFromAttributeMap(target map[string]any, attr_source_map map[string]attr.Value, names ...string) {
	for _, name := range names {
		if attr_source_map[name] != nil && !attr_source_map[name].IsNull() {
			target[name] = GetAttributeValue[String](attr_source_map, name).ValueString()
		}
	}
}

// Sets one or more optional string attributes from a source map
func SetOptionalAttributeStringFromMap(target map[string]attr.Value, source map[string]any, names ...string) {
	for _, name := range names {
		if source[name] != nil {
			target[name] = StringValue(source[name].(string))
		}
	}
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

// Kafka-specific functions, shared between multiple types
func BrokersFromModelList(Brokers List, dd diag.Diagnostics) ([]map[string]any, diag.Diagnostics) {
	output := make([]map[string]any, 0)
	elements := Brokers.Elements()
	for _, b := range elements {
		broker := map[string]any{}
		attrs := b.(basetypes.ObjectValue).Attributes()
		for k, v := range attrs {
			switch v.(type) {
			case String:
				value, ok := attrs[k].(basetypes.StringValue)
				if !ok {
					dd.AddError(
						"Could not look up attribute value",
						fmt.Sprintf("Cannot cast key %s to a string value. Please report this to Mezmo.", k),
					)
					continue
				}
				broker[k] = value.ValueString()
			case Int64:
				value, ok := attrs[k].(basetypes.Int64Value)
				if !ok {
					dd.AddError(
						"Could not look up attribute value",
						fmt.Sprintf("Cannot cast key %s to an int value. Please report this to Mezmo.", k),
					)
					continue
				}
				broker[k] = value.ValueInt64()
			}
		}
		output = append(output, broker)
	}
	return output, dd
}

func BrokersToModelList(elementType attr.Type, brokers []interface{}) List {
	output := make([]attr.Value, 0)
	for _, v := range brokers {
		broker_raw := v.(map[string]interface{})
		broker_map := map[string]attr.Value{
			"host": StringValue(broker_raw["host"].(string)),
			"port": Int64Value(int64(broker_raw["port"].(float64))),
		}
		broker := basetypes.NewObjectValueMust(
			map[string]attr.Type{"host": StringType, "port": Int64Type},
			broker_map)
		output = append(output, broker)
	}
	brokersList, _ := ListValue(elementType, output)
	return brokersList
}

func KafkaDestinationSASLToModel(types map[string]attr.Type, user_config map[string]interface{}) basetypes.ObjectValue {
	sasl := map[string]attr.Value{}
	if user_config["sasl_username"] != nil {
		sasl["username"] = StringValue(user_config["sasl_username"].(string))
	}
	if user_config["sasl_password"] != nil {
		sasl["password"] = StringValue(user_config["sasl_password"].(string))
	}
	if user_config["sasl_mechanism"] != nil {
		sasl["mechanism"] = StringValue(user_config["sasl_mechanism"].(string))
	}
	return basetypes.NewObjectValueMust(types, sasl)
}
