package modelutils

import (
	"reflect"

	"fmt"

	"slices"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// This function can help with taking an object from an API response and translating
// it into a basetype object from the terraform-plugin-framework
func MapAnyToMapValues(values map[string]any) map[string]attr.Value {
	result := make(map[string]attr.Value, len(values))
	for k, v := range values {
		// Add other types here if you need
		switch v.(type) {
		case string:
			result[k] = StringValue(v.(string))
		default:
			panic(fmt.Errorf("unsupported value type %T for key %v", v, k))
		}
	}
	return result
}

// Make a complete map of values, using a null-equivalent for any missing fields
func MapAnyFillMissingValues(attrs map[string]attr.Type, values map[string]any, optional_fields []string) map[string]attr.Value {
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
				result[k] = SliceToStringListValue(v.([]any))
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
// return a regular map[string]any object that can be used in `Component` such that
// it can be used in an API call.
func MapValuesToMapAny(obj interface{}, dd *diag.Diagnostics) map[string]any {
	var attrs map[string]attr.Value
	switch obj := obj.(type) {
	case Object:
		attrs = obj.Attributes()
	case Map:
		attrs = obj.Elements()
	default:
		panic(fmt.Sprintf("Unsupported object type for `fromAttributes`: %T", obj))
	}

	target := make(map[string]any, len(attrs))

	for k, v := range attrs {
		if v.IsUnknown() || v.IsNull() {
			continue
		}
		switch v := v.(type) {
		// Add more types here if need be
		case Bool:
			target[k] = v.ValueBool()
		case String:
			target[k] = v.ValueString()
		case Int64:
			target[k] = v.ValueInt64()
		case List:
			target[k] = StringListValueToStringSlice(v)
		default:
			panic("unhandled type (add new type if it needs to be handled)")
		}
	}
	return target
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

// Given a map of schema (resource, data, etc) attributes, returns the attribute types
func ToAttrTypes(attributes map[string]schema.Attribute) map[string]attr.Type {
	result := make(map[string]attr.Type)
	for k, v := range attributes {
		result[k] = v.GetType()
	}

	return result
}
