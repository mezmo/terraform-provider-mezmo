package modelutils

const VRL_PARSER_APACHE = "apache_log"
const VRL_PARSER_CEF = "cef_log"
const VRL_PARSER_CSV = "csv_row"
const VRL_PARSER_GROK = "grok_parser"
const VRL_PARSER_KEY_VALUE = "key_value_log"
const VRL_PARSER_NGINX = "nginx_log"
const VRL_PARSER_REGEX = "regex_parser"
const VRL_PARSER_TIMESTAMP = "timestamp_parser"

var VRL_PARSERS = map[string]string{
	"apache_log":              "parse_apache_log",
	"aws_alb_log":             "parse_aws_alb_log",
	"aws_cloudwatch_log":      "parse_aws_cloudwatch_log_subscription_message",
	"aws_vpc_flow_log":        "parse_aws_vpc_flow_log",
	"cef_log":                 "parse_cef",
	"common_log":              "parse_common_log",
	"csv_row":                 "parse_csv",
	"glog":                    "parse_glog",
	"grok_parser":             "parse_grok",
	"int_parser":              "parse_int",
	"json_parser":             "parse_json",
	"key_value_log":           "parse_key_value",
	"klog":                    "parse_klog",
	"linux_authorization_log": "parse_linux_authorization",
	"nginx_log":               "parse_nginx_log",
	"querystring_parser":      "parse_query_string",
	"regex_parser":            "parse_regex",
	"syslog":                  "parse_syslog",
	"timestamp_parser":        "parse_timestamp",
	"token_parser":            "parse_tokens",
	"url_parser":              "parse_url",
	"user_agent_parser":       "parse_user_agent",
}

var VRL_PARSERS_WITH_REQUIRED_OPTIONS = []string{
	"apache_log",
	"grok_parser",
	"regex_parser",
	"nginx_log",
	"timestamp_parser",
}

var OPTIONAL_FIELDS_BY_PARSER = map[string][]string{
	"apache_log":       {"timestamp_format", "custom_timestamp_format"},
	"csv_row":          {"field_names"},
	"key_value_log":    {"field_delimiter", "key_delimiter"},
	"nginx_log":        {"timestamp_format", "custom_timestamp_format"},
	"timestamp_parser": {"custom_format"},
}

var NGINX_LOG_FORMATS = []string{"combined", "error"}
var APACHE_LOG_FORMATS = append(NGINX_LOG_FORMATS, "common")
