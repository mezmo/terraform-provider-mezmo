package modelutils

var Non_Change_Operator_Labels = []string{
	"contains",
	"ends_with",
	"equal",
	"exists",
	"greater",
	"greater_or_equal",
	"is_array",
	"is_boolean",
	"is_empty",
	"is_ip_in_cidr_range",
	"is_metric",
	"is_null",
	"is_number",
	"is_object",
	"is_string",
	"less",
	"less_or_equal",
	"regex_match",
	"starts_with",
}

var Change_Operator_Labels = []string{
	"percent_change_greater",
	"percent_change_greater_or_equal",
	"percent_change_less",
	"percent_change_less_or_equal",
	"value_change_greater",
	"value_change_greater_or_equal",
	"value_change_less",
	"value_change_less_or_equal",
}
