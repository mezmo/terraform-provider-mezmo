package processors

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
)

var base_options_parser_schema = SchemaAttributes{
	"parser": schema.StringAttribute{
		Required:    true,
		Description: "The kind of parser to use against the input value from \"field\".",
		Validators: []validator.String{
			stringvalidator.OneOf(modelutils.MapKeys(modelutils.VRL_PARSERS)...),
		},
	},
	"apache_log_options": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "Options for apache log parser",
		Attributes: map[string]schema.Attribute{
			"format": schema.StringAttribute{
				Required:    true,
				Description: "The log format.",
				Validators: []validator.String{
					stringvalidator.OneOf(modelutils.APACHE_LOG_FORMATS...),
				},
			},
			"timestamp_format": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The timestamp format of log entries.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"custom_timestamp_format": schema.StringAttribute{
				Optional:    true,
				Description: "Custom timestamp format, according to strftime, for log entries.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	},
	"cef_log_options": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "Options for common event format (CEF) log parser",
		Attributes: map[string]schema.Attribute{
			"translate_custom_fields": schema.BoolAttribute{
				Optional:    true,
				Description: "Translate custom fields in log.",
			},
		},
	},
	"csv_row_options": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "Options for CSV row parser",
		Attributes: map[string]schema.Attribute{
			"field_names": schema.ListAttribute{
				ElementType: StringType,
				Optional:    true,
				Description: "The name of the CSV fields that take the value in the same order " +
					"they appear in data",
				Validators: []validator.List{
					listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				},
			},
		},
	},
	"grok_parser_options": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "Options for grok parser",
		Attributes: map[string]schema.Attribute{
			"pattern": schema.StringAttribute{
				Required:    true,
				Description: "The grok pattern. Must be composed of community expressions.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	},
	"key_value_log_options": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "Options for key value log parser",
		Attributes: map[string]schema.Attribute{
			"key_value_delimiter": schema.StringAttribute{
				Optional:    true,
				Description: "One or more characters that separate each key and value.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"field_delimiter": schema.StringAttribute{
				Optional:    true,
				Description: "One or more characters that separate each key/value pair.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	},
	"nginx_log_options": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "Options for nginx log parser",
		Attributes: map[string]schema.Attribute{
			"format": schema.StringAttribute{
				Required:    true,
				Description: "The log format.",
				Validators: []validator.String{
					stringvalidator.OneOf(modelutils.NGINX_LOG_FORMATS...),
				},
			},
			"timestamp_format": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The timestamp format of log entries.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"custom_timestamp_format": schema.StringAttribute{
				Optional:    true,
				Description: "Custom timestamp format, according to strftime, for log entries.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	},
	"regex_parser_options": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "Options for regex parser",
		Attributes: map[string]schema.Attribute{
			"pattern": schema.StringAttribute{
				Required:    true,
				Description: "The regex pattern.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"case_sensitive": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Perform case sensitive regex matching",
			},
			"multiline": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "When true, ^ and $ will match not only the beginning and end of a string but also the start and end of lines",
			},
			"match_newline": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "When true, dot (.) will match newline characters",
			},
			"crlf_newline": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "When true and multiline = true, newline character will be interpreted as \r\n",
			},
		},
	},
	"timestamp_parser_options": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "Options for timestamp log parser",
		Attributes: map[string]schema.Attribute{
			"format": schema.StringAttribute{
				Required:    true,
				Description: "The timestamp format.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"custom_format": schema.StringAttribute{
				Optional:    true,
				Description: "Custom timestamp format, according to strftime, for log entries.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	},
}
var base_schema = SchemaAttributes{
	"field": schema.StringAttribute{
		Required:    true,
		Description: "The JSON field whose value should be parsed.",
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	},
	"target_field": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Description: "The field into which the parsed value should be inserted. Leave blank to " +
			"insert the parsed data into the original field.",
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	},
}

var parse_schema = ExtendSchemaAttributes(base_schema, copySchema(base_options_parser_schema))
var parse_sequential_item_schema = ExtendSchemaAttributes(SchemaAttributes{
	"label": schema.StringAttribute{
		Optional:    true,
		Description: "An arbitrary label to identify the parser",
	},
	"output_name": schema.StringAttribute{
		Computed: true,
		Description: "A system generated value to identify the results of this parser. " +
			"This value should be used when connecting the results to another processor or destination.",
	},
}, copySchema(base_options_parser_schema))

var parse_sequential_schema = ExtendSchemaAttributes(base_schema, SchemaAttributes{
	"parsers": schema.ListNestedAttribute{
		Required:    true,
		Description: "The ordered list of parsers to use to parse the field. Parsing short-circuits on the first successful parse attempt.",
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: parse_sequential_item_schema,
		},
	},
})

// Make copy of the parse base schema to ensure it can be shared
// by single and sequential parser
func copySchema(base_schema SchemaAttributes) SchemaAttributes {
	new_schema := make(map[string]schema.Attribute)
	for k, v := range base_schema {
		new_schema[k] = v
	}
	return new_schema
}

// Given a parser, generates the parser options key. Example, apache_log -> apache_log_options
func optionsKeyForParser(parser string) string {
	return fmt.Sprintf("%s_options", parser)
}

// Given a parser options key, such as apache_log_options, returns the schema attributes
func optionsSchemaAttributes(options_key string) (map[string]schema.Attribute, error) {
	options_schema, ok := base_options_parser_schema[options_key]

	if !ok {
		return nil, fmt.Errorf("unknown parser option: %s", options_key)
	}

	return options_schema.(schema.SingleNestedAttribute).Attributes, nil
}
