package processors

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/providertest"
)

func TestAccParseProcessor(t *testing.T) {
	const cacheKey = "parse_resources"
	SetCachedConfig(cacheKey, `
		resource "mezmo_pipeline" "test_parent" {
			title = "pipeline"
		}
		resource "mezmo_http_source" "my_source" {
			pipeline_id = mezmo_pipeline.test_parent.id
		}`,
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Error: pipeline_id is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						parser = "parse_apache_log"
						apache_log_options = {
							format = "error"
						}
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},

			// Error: `field` is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						parser = "apache_log"
						apache_log_options = {
							format = "common"
						}
					}`,
				ExpectError: regexp.MustCompile("The argument \"field\" is required"),
			},

			// Error: `parser` is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".field"
						apache_log_options = {
							format = "common"
						}
					}`,
				ExpectError: regexp.MustCompile("The argument \"parser\" is required"),
			},

			// Error: `parser` is an enum
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".field"
						parser = "my_own_parser"
						apache_log_options = {
							format = "common"
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute parser value must be one of:"),
			},

			// apache_log_options is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".field"
						parser = "apache_log"
					}`,
				ExpectError: regexp.MustCompile("Attribute apache_log_options is required for apache_log."),
			},
			// apache format is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".field"
						parser = "parse_apache_log"
						apache_log_options = {
							format = ""
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute apache_log_options.format value must be one of"),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".field"
						parser = "parse_apache_log"
						apache_log_options = {
							timestamp_format = "Custom"
						}
					}`,
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"apache_log_options\": attribute \"format\""),
			},
			// apache format must be a known value
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".field"
						parser = "parse_apache_log"
						apache_log_options = {
							format = "unknown"
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute apache_log_options.format value must be one of"),
			},
			// apache timestamp format, if specified, cannot be empty
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".field"
						parser = "parse_apache_log"
						apache_log_options = {
							format = "common"
							timestamp_format = ""
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute apache_log_options.timestamp_format string length must be at least"),
			},
			// apache custom timestamp format, if specified, cannot be empty
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".field"
						parser = "parse_apache_log"
						apache_log_options = {
							format = "common"
							timestamp_format = "Custom"
							custom_timestamp_format = ""
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute apache_log_options.custom_timestamp_format string length must be"),
			},

			// nginx_log_options is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".field"
						parser = "nginx_log"
					}`,
				ExpectError: regexp.MustCompile("Attribute nginx_log_options is required for nginx_log."),
			},
			// Error: nginx_format must be one of the allowed values
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".field"
						parser = "nginx_log"
						nginx_log_options = {
							format = "unknown"
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute nginx_log_options.format value must be one of"),
			},
			// Error: nginx format is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".field"
						parser = "nginx_log"
						nginx_log_options = {
						}
					}`,
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"nginx_log_options\": attribute \"format\""),
			},
			// nginx timestamp format, if specified, cannot be empty
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".field"
						parser = "nginx_log"
						nginx_log_options = {
							format = "common"
							timestamp_format = ""
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute nginx_log_options.timestamp_format string length must be at least"),
			},
			// nginx custom timestamp format, if specified, cannot be empty
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".field"
						parser = "nginx_log"
						nginx_log_options = {
							format = "common"
							timestamp_format = "Custom"
							custom_timestamp_format = ""
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute nginx_log_options.custom_timestamp_format string length must be"),
			},

			// grok_parser_options is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".field"
						parser = "grok_parser"
					}`,
				ExpectError: regexp.MustCompile("Attribute grok_parser_options is required for grok_parser."),
			},
			// Error: pattern is required for grok parser
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".field"
						parser = "grok_parser"
						grok_parser_options = {}
					}`,
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"grok_parser_options\": attribute \"pattern\""),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".field"
						parser = "grok_parser"
						grok_parser_options = {
							pattern = ""
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute grok_parser_options.pattern string length must be at least 1"),
			},

			// regex_parser_options is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".field"
						parser = "regex_parser"
					}`,
				ExpectError: regexp.MustCompile("Attribute regex_parser_options is required for regex_parser."),
			},
			// Error: pattern is required for regex parser
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".field"
						parser = "regex_parser"
						regex_parser_options = {}
					}`,
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"regex_parser_options\": attribute \"pattern\""),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".field"
						parser = "regex_parser"
						regex_parser_options = {
							pattern = ""
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute regex_parser_options.pattern string length must be at least 1"),
			},

			// timestamp_parser_options is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".field"
						parser = "timestamp_parser"
					}`,
				ExpectError: regexp.MustCompile("Attribute timestamp_parser_options is required for timestamp_parser."),
			},
			// Error: format is required for timestamp parser
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".field"
						parser = "timestamp_parser"
						timestamp_parser_options = {}
					}`,
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"timestamp_parser_options\": attribute"),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".field"
						parser = "timestamp_parser"
						timestamp_parser_options = {
							format = ""
						}
					}`,
				// min length = 1
				ExpectError: regexp.MustCompile("Attribute timestamp_parser_options.format string length must be"),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".field"
						parser = "timestamp_parser"
						timestamp_parser_options = {
							format = "Custom"
							custom_format = ""
						}
					}`,
				// length at least 1
				ExpectError: regexp.MustCompile("Attribute timestamp_parser_options.custom_format string length must be"),
			},

			// apache log - default timestamp
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "apache_default" {
						title = "custom apache parser title"
						description = "custom apache parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parser = "apache_log"
						apache_log_options = {
							format = "combined"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_processor.apache_default", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_processor.apache_default", map[string]any{
						"pipeline_id":                         "#mezmo_pipeline.test_parent.id",
						"title":                               "custom apache parser title",
						"description":                         "custom apache parser desc",
						"generation_id":                       "0",
						"inputs.#":                            "0",
						"field":                               ".something",
						"parser":                              "apache_log",
						"apache_log_options.format":           "combined",
						"apache_log_options.timestamp_format": "%d/%b/%Y:%T %z",
					}),
				),
			},
			// apache log - custom timestamp
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "apache_default_custom_time" {
						title = "custom apache parser title"
						description = "custom apache parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parser = "apache_log"
						apache_log_options = {
							format = "combined"
							timestamp_format = "Custom"
							custom_timestamp_format = "%d-%b-%YT%T"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_processor.apache_default_custom_time", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_processor.apache_default_custom_time", map[string]any{
						"pipeline_id":                         "#mezmo_pipeline.test_parent.id",
						"title":                               "custom apache parser title",
						"description":                         "custom apache parser desc",
						"generation_id":                       "0",
						"inputs.#":                            "0",
						"field":                               ".something",
						"parser":                              "apache_log",
						"apache_log_options.format":           "combined",
						"apache_log_options.timestamp_format": "Custom",
						"apache_log_options.custom_timestamp_format": "%d-%b-%YT%T",
					}),
				),
			},
			// cef log - no options
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "cef_default" {
						title = "cef parser title"
						description = "cef parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parser = "cef_log"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_processor.cef_default", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_processor.cef_default", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "cef parser title",
						"description":   "cef parser desc",
						"generation_id": "0",
						"inputs.#":      "0",
						"field":         ".something",
						"parser":        "cef_log",
					}),
				),
			},
			// cef log with options
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "cef_translate_custom" {
						title = "cef parser title"
						description = "cef parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parser = "cef_log"
						cef_log_options = {
							translate_custom_fields = true
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_processor.cef_translate_custom", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_processor.cef_translate_custom", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "cef parser title",
						"description":   "cef parser desc",
						"generation_id": "0",
						"inputs.#":      "0",
						"field":         ".something",
						"parser":        "cef_log",
						"cef_log_options.translate_custom_fields": "true",
					}),
				),
			},
			// CSV parser - no options
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "csv_parser_default" {
						title = "custom csv parser title"
						description = "custom csv parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parser = "csv_row"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_processor.csv_parser_default", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_processor.csv_parser_default", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "custom csv parser title",
						"description":   "custom csv parser desc",
						"generation_id": "0",
						"inputs.#":      "0",
						"field":         ".something",
						"parser":        "csv_row",
					}),
				),
			},
			// CSV parser with options
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "csv_parser" {
						title = "custom csv parser title"
						description = "custom csv parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parser = "csv_row"
						csv_row_options = {
							field_names = ["field1", "field2"]
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_processor.csv_parser", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_processor.csv_parser", map[string]any{
						"pipeline_id":                   "#mezmo_pipeline.test_parent.id",
						"title":                         "custom csv parser title",
						"description":                   "custom csv parser desc",
						"generation_id":                 "0",
						"inputs.#":                      "0",
						"field":                         ".something",
						"parser":                        "csv_row",
						"csv_row_options.field_names.#": "2",
						"csv_row_options.field_names.0": "field1",
						"csv_row_options.field_names.1": "field2",
					}),
				),
			},
			// grok parser
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "grok_parser" {
						title = "grok parser title"
						description = "grok parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parser = "grok_parser"
						grok_parser_options = {
							pattern = "%%{USERNAME:user} %%{EMAILADDRESS:email}"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_processor.grok_parser", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_processor.grok_parser", map[string]any{
						"pipeline_id":                 "#mezmo_pipeline.test_parent.id",
						"title":                       "grok parser title",
						"description":                 "grok parser desc",
						"generation_id":               "0",
						"inputs.#":                    "0",
						"field":                       ".something",
						"parser":                      "grok_parser",
						"grok_parser_options.pattern": "%{USERNAME:user} %{EMAILADDRESS:email}",
					}),
				),
			},
			// key value parser - no options
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "key_value_parser" {
						title = "key value parser title"
						description = "key value parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parser = "key_value_log"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_processor.key_value_parser", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_processor.key_value_parser", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "key value parser title",
						"description":   "key value parser desc",
						"generation_id": "0",
						"inputs.#":      "0",
						"field":         ".something",
						"parser":        "key_value_log",
					}),
				),
			},
			// key value with all options
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "key_value_parser_all_options" {
						title = "key value parser title"
						description = "key value parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parser = "key_value_log"
						key_value_log_options = {
							field_delimiter = ":"
							key_value_delimiter = ","
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_processor.key_value_parser_all_options", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_processor.key_value_parser_all_options", map[string]any{
						"pipeline_id":                           "#mezmo_pipeline.test_parent.id",
						"title":                                 "key value parser title",
						"description":                           "key value parser desc",
						"generation_id":                         "0",
						"inputs.#":                              "0",
						"field":                                 ".something",
						"parser":                                "key_value_log",
						"key_value_log_options.field_delimiter": ":",
						"key_value_log_options.key_value_delimiter": ",",
					}),
				),
			},
			// field delimiter only
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "key_value_parser_field_only" {
						title = "key value parser title"
						description = "key value parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parser = "key_value_log"
						key_value_log_options = {
							field_delimiter = ":"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_processor.key_value_parser_field_only", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_processor.key_value_parser_field_only", map[string]any{
						"pipeline_id":                           "#mezmo_pipeline.test_parent.id",
						"title":                                 "key value parser title",
						"description":                           "key value parser desc",
						"generation_id":                         "0",
						"inputs.#":                              "0",
						"field":                                 ".something",
						"parser":                                "key_value_log",
						"key_value_log_options.field_delimiter": ":",
					}),
					StateDoesNotHaveFields("mezmo_parse_processor.key_value_parser_field_only", []string{"key_value_log_options.key_value_delimiter"}),
				),
			},
			// key value field delimiter only
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "key_value_parser_key_only" {
						title = "key value parser title"
						description = "key value parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parser = "key_value_log"
						key_value_log_options = {
							key_value_delimiter = ":"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_processor.key_value_parser_key_only", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_processor.key_value_parser_key_only", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "key value parser title",
						"description":   "key value parser desc",
						"generation_id": "0",
						"inputs.#":      "0",
						"field":         ".something",
						"parser":        "key_value_log",
						"key_value_log_options.key_value_delimiter": ":",
					}),
					StateDoesNotHaveFields("mezmo_parse_processor.key_value_parser_key_only", []string{"key_value_log_options.field_delimiter"}),
				),
			},
			// nginx default
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "nginx_default" {
						title = "nginx parser title"
						description = "nginx parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parser = "nginx_log"
						nginx_log_options = {
							format = "error"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_processor.nginx_default", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_processor.nginx_default", map[string]any{
						"pipeline_id":                        "#mezmo_pipeline.test_parent.id",
						"title":                              "nginx parser title",
						"description":                        "nginx parser desc",
						"generation_id":                      "0",
						"inputs.#":                           "0",
						"field":                              ".something",
						"parser":                             "nginx_log",
						"nginx_log_options.format":           "error",
						"nginx_log_options.timestamp_format": "%Y/%m/%d %H:%M:%S",
					}),
					StateDoesNotHaveFields("mezmo_parse_processor.nginx_default", []string{"nginx_log_options.custom_timestamp_format"}),
				),
			},
			// nginx custom timestamp
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "nginx_custom_time" {
						title = "nginx parser title"
						description = "nginx parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parser = "nginx_log"
						nginx_log_options = {
							format = "error"
							timestamp_format = "Custom"
							custom_timestamp_format = "%d-%b-%YT%T"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_processor.nginx_custom_time", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_processor.nginx_custom_time", map[string]any{
						"pipeline_id":                        "#mezmo_pipeline.test_parent.id",
						"title":                              "nginx parser title",
						"description":                        "nginx parser desc",
						"generation_id":                      "0",
						"inputs.#":                           "0",
						"field":                              ".something",
						"parser":                             "nginx_log",
						"nginx_log_options.format":           "error",
						"nginx_log_options.timestamp_format": "Custom",
						"nginx_log_options.custom_timestamp_format": "%d-%b-%YT%T",
					}),
				),
			},
			// regex parser without flags
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "regex_defaults" {
						title = "regex parser title"
						description = "regex parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parser = "regex_parser"
						regex_parser_options = {
							pattern = "\\d{3}-\\d{2}-\\d{3}"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_processor.regex_defaults", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_processor.regex_defaults", map[string]any{
						"pipeline_id":                         "#mezmo_pipeline.test_parent.id",
						"title":                               "regex parser title",
						"description":                         "regex parser desc",
						"generation_id":                       "0",
						"inputs.#":                            "0",
						"field":                               ".something",
						"parser":                              "regex_parser",
						"regex_parser_options.pattern":        "\\d{3}-\\d{2}-\\d{3}",
						"regex_parser_options.case_sensitive": "true",
						"regex_parser_options.multiline":      "false",
						"regex_parser_options.match_newline":  "false",
						"regex_parser_options.crlf_newline":   "false",
					}),
				),
			},
			// regex parser with flags
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "regex_custom" {
						title = "regex parser title"
						description = "regex parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parser = "regex_parser"
						regex_parser_options = {
							pattern = "\\d{3}-\\d{2}-\\d{3}"
							case_sensitive = false
							multiline = true
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_processor.regex_custom", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_processor.regex_custom", map[string]any{
						"pipeline_id":                         "#mezmo_pipeline.test_parent.id",
						"title":                               "regex parser title",
						"description":                         "regex parser desc",
						"generation_id":                       "0",
						"inputs.#":                            "0",
						"field":                               ".something",
						"parser":                              "regex_parser",
						"regex_parser_options.pattern":        "\\d{3}-\\d{2}-\\d{3}",
						"regex_parser_options.case_sensitive": "false",
						"regex_parser_options.multiline":      "true",
						"regex_parser_options.match_newline":  "false",
						"regex_parser_options.crlf_newline":   "false",
					}),
				),
			},
			// timestamp parser default
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "timestamp_parser_default" {
						title = "timestamp parser title"
						description = "timestamp parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parser = "timestamp_parser"
						timestamp_parser_options = {
							format = "%a %d %b %T %Y"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_processor.timestamp_parser_default", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_processor.timestamp_parser_default", map[string]any{
						"pipeline_id":                     "#mezmo_pipeline.test_parent.id",
						"title":                           "timestamp parser title",
						"description":                     "timestamp parser desc",
						"generation_id":                   "0",
						"inputs.#":                        "0",
						"field":                           ".something",
						"parser":                          "timestamp_parser",
						"timestamp_parser_options.format": "%a %d %b %T %Y",
					}),
					StateDoesNotHaveFields("mezmo_parse_processor.timestamp_parser_default", []string{"timestamp_parser_options.custom_format"}),
				),
			},
			// timestamp parser custom format
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "timestamp_parser_custom" {
						title = "timestamp parser title"
						description = "timestamp parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parser = "timestamp_parser"
						timestamp_parser_options = {
							format = "Custom"
							custom_format = "%Y/%m/%d %H:%M:%S"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_processor.timestamp_parser_custom", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_processor.timestamp_parser_custom", map[string]any{
						"pipeline_id":                            "#mezmo_pipeline.test_parent.id",
						"title":                                  "timestamp parser title",
						"description":                            "timestamp parser desc",
						"generation_id":                          "0",
						"inputs.#":                               "0",
						"field":                                  ".something",
						"parser":                                 "timestamp_parser",
						"timestamp_parser_options.format":        "Custom",
						"timestamp_parser_options.custom_format": "%Y/%m/%d %H:%M:%S",
					}),
				),
			},
			// Parser with empty target_field
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "with_empty_target" {
						title = "custom apache parser title"
						description = "custom apache parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						target_field = ""
						parser = "apache_log"
						apache_log_options = {
							format = "common"
							timestamp_format = "Custom"
							custom_timestamp_format = "%Y/%m/%d %H:%M:%S"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_processor.with_empty_target", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_processor.with_empty_target", map[string]any{
						"pipeline_id":                         "#mezmo_pipeline.test_parent.id",
						"title":                               "custom apache parser title",
						"description":                         "custom apache parser desc",
						"generation_id":                       "0",
						"inputs.#":                            "0",
						"field":                               ".something",
						"target_field":                        "",
						"parser":                              "apache_log",
						"apache_log_options.format":           "common",
						"apache_log_options.timestamp_format": "Custom",
						"apache_log_options.custom_timestamp_format": "%Y/%m/%d %H:%M:%S",
					}),
				),
			},
			// with non empty parser
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "with_target" {
						title = "custom apache parser title"
						description = "custom apache parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						target_field = ".parsed"
						parser = "apache_log"
						apache_log_options = {
							format = "common"
							timestamp_format = "Custom"
							custom_timestamp_format = "%Y/%m/%d %H:%M:%S"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_processor.with_target", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_processor.with_target", map[string]any{
						"pipeline_id":                         "#mezmo_pipeline.test_parent.id",
						"title":                               "custom apache parser title",
						"description":                         "custom apache parser desc",
						"generation_id":                       "0",
						"inputs.#":                            "0",
						"field":                               ".something",
						"target_field":                        ".parsed",
						"parser":                              "apache_log",
						"apache_log_options.format":           "common",
						"apache_log_options.timestamp_format": "Custom",
						"apache_log_options.custom_timestamp_format": "%Y/%m/%d %H:%M:%S",
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "import_target" {
						title = "custom apache parser title"
						description = "custom apache parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parser = "apache_log"
						apache_log_options = {
							format = "common"
							timestamp_format = "Custom"
							custom_timestamp_format = "%Y/%m/%d %H:%M:%S"
						}
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_parse_processor.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_parse_processor.with_target"),
				ImportStateVerify: true,
			},

			// Update apache parser with custom timestamp to nginx
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "with_target" {
						title = "custom apache parser title"
						description = "custom apache parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						field = ".something"
						parser = "nginx_log"
						nginx_log_options = {
							format = "combined"
							timestamp_format = "Custom"
							custom_timestamp_format = "%a %b %e %T %Y"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_processor.with_target", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_processor.with_target", map[string]any{
						"pipeline_id":                        "#mezmo_pipeline.test_parent.id",
						"title":                              "custom apache parser title",
						"description":                        "custom apache parser desc",
						"generation_id":                      "1",
						"inputs.#":                           "1",
						"inputs.0":                           "#mezmo_http_source.my_source.id",
						"field":                              ".something",
						"target_field":                       "",
						"parser":                             "nginx_log",
						"nginx_log_options.format":           "combined",
						"nginx_log_options.timestamp_format": "Custom",
						"nginx_log_options.custom_timestamp_format": "%a %b %e %T %Y",
					}),
				),
			},
			// // Server side validation on create
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "apache_defaults" {
						title = "apache parser title"
						description = "apache parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parser = "apache_log"
						apache_log_options = {
							format = "common"
							timestamp_format = "Custom"
						}
					}`,
				// custom_timestamp_format field is required
				ExpectError: regexp.MustCompile("(?s)have required property.*'custom_timestamp_format'"),
			},

			// // Error: server-side validation on update
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_processor" "with_target" {
						title = "custom apache parser title"
						description = "custom apache parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						field = ".something"
						parser = "nginx_log"
						nginx_log_options = {
							format = "combined"
							timestamp_format = "Custom"
						}
					}`,
				// custom_timestamp_format field is required
				ExpectError: regexp.MustCompile("(?s)have required property.*'custom_timestamp_format'"),
			},
			// confirm manually deleted resources are recreated
			{
				Config: GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_parse_processor" "test_processor" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					inputs 			= []
					field				= ".field"
					parser 			= "apache_log"
					apache_log_options = {
						format = "combined"
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_processor.test_processor", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_parse_processor.test_processor", "title", "new title"),
					// delete the resource
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_parse_processor.test_processor",
					),
				),
				// verify resource will be re-created after refresh
				ExpectNonEmptyPlan: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
