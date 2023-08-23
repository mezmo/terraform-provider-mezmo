package modelutils

var Operators = []string{
	"greater", "greater_or_equal", "less", "less_or_equal", "equal", "not_equal",
	"contains", "is_ip_in_cidr_range",
	"is_metric", "is_array", "is_boolean", "is_empty", "is_null", "is_number", "is_object", "is_string",
	"exists", "not_exists",
	"ends_with", "starts_with", "regex_match", "not_regex_match",
}
