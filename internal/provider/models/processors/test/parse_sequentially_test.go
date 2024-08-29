package processors

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/providertest"
)

func TestAccParseSequentiallyProcessor(t *testing.T) {
	const cacheKey = "parse_sequentially_resources"
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
					resource "mezmo_parse_sequentially_processor" "my_processor" {
						field = ".app"
						parsers = [
							{
								parser = "parse_apache_log"
								apache_log_options = {
									format = "error"
								}
							}
						]
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},
			// Error: field is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						parsers = [
							{
								parser = "parse_apache_log"
								apache_log_options = {
									format = "error"
								}
							}
						]
					}`,
				ExpectError: regexp.MustCompile("The argument \"field\" is required"),
			},

			// Error: `parsers` is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
					}`,
				ExpectError: regexp.MustCompile("(?s)The argument \"parsers\" is required, but no definition .*was found"),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
						parsers = []
					}`,
				ExpectError: regexp.MustCompile("(?s)Attribute parsers list must contain at least 1 elements, got:.* 0"),
			},

			// Error: `parser` is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
						parsers = [
							{
								label = "apache errors"
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`(?s)Inappropriate value for attribute "parsers": element 0: attribute "parser" is.*required.`),
			},
			// parser is an enum
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
						parsers = [
							{
								label = "apache errors"
								parser = "unknown"
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`Attribute parsers\[0\].parser value must be one of:`),
			},

			// apache_log_options is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
						parsers = [
							{
								label = "apache errors"
								parser = "apache_log"
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`(?s)Inappropriate value for attribute "parsers": element 0: attribute.*"apache_log_options" is required.`),
			},
			// apache_log_options format must be an enum
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
						parsers = [
							{
								label = "apache errors"
								parser = "apache_log"
								apache_log_options = {
									format = "unknown"
								}
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`Attribute parsers\[0\].apache_log_options.format value must be one of:`),
			},
			// apache_log_options format is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
						parsers = [
							{
								label = "apache errors"
								parser = "apache_log"
								apache_log_options = {
								}
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`(?s)Inappropriate value for attribute "parsers": element 0: attribute.*`),
			},
			// apache_log_options timestamp_format cannot be empty
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
						parsers = [
							{
								label = "apache errors"
								parser = "apache_log"
								apache_log_options = {
									format = "common"
									timestamp_format = ""
								}
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`(?s)Attribute parsers\[0\].apache_log_options.timestamp_format string length must.*be at least 1, got: 0`),
			},
			// apache_log_options custom_timestamp_format cannot be empty
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
						parsers = [
							{
								label = "apache errors"
								parser = "apache_log"
								apache_log_options = {
									format = "common"
									timestamp_format = "Custom"
									custom_timestamp_format = ""
								}
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`(?s)Attribute parsers\[0\].apache_log_options.custom_timestamp_format string length.*must be at least 1, got: 0`),
			},

			// nginx_log_options is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "nginx_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
						parsers = [
							{
								label = "error logs"
								parser = "nginx_log"
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`(?s)Inappropriate value for attribute "parsers": element 0: attribute.*"nginx_log_options" is required.`),
			},
			// nginx_log_options format must be an enum
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "nginx_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
						parsers = [
							{
								label = "nginx errors"
								parser = "nginx_log"
								nginx_log_options = {
									format = "unknown"
								}
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`Attribute parsers\[0\].nginx_log_options.format value must be one of:`),
			},
			// nginx_log_options format is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "nginx_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
						parsers = [
							{
								label = "nginx errors"
								parser = "nginx_log"
								nginx_log_options = {
								}
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`(?s)Inappropriate value for attribute "parsers": element 0: attribute.*`),
			},
			// nginx_log_options timestamp_format cannot be empty
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "nginx_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
						parsers = [
							{
								label = "nginx errors"
								parser = "nginx_log"
								nginx_log_options = {
									format = "combined"
									timestamp_format = ""
								}
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`(?s)Attribute parsers\[0\].nginx_log_options.timestamp_format string length must be.*at least 1, got: 0`),
			},
			// nginx_log_options custom_timestamp_format cannot be empty
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "nginx_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
						parsers = [
							{
								label = "error logs"
								parser = "nginx_log"
								nginx_log_options = {
									format = "common"
									timestamp_format = "Custom"
									custom_timestamp_format = ""
								}
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`(?s)Attribute parsers\[0\].nginx_log_options.custom_timestamp_format string length.*must be at least 1, got: 0`),
			},

			// grok_parser
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
						parsers = [
							{
								label = "error logs"
								parser = "grok_parser"
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`(?s)Inappropriate value for attribute "parsers": element 0: attribute.*"grok_parser_options" is required.`),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
						parsers = [
							{
								label = "error logs"
								parser = "grok_parser"
								grok_parser_options = {
								}
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`(?s)Inappropriate value for attribute "parsers": element 0: attribute.*"grok_parser_options": attribute "pattern" is required.`),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
						parsers = [
							{
								label = "error logs"
								parser = "grok_parser"
								grok_parser_options = {
									pattern = ""
								}
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`(?s)Attribute parsers\[0\].grok_parser_options.pattern string length must be at.*least 1, got: 0`),
			},

			// regex_parser
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
						parsers = [
							{
								label = "error logs"
								parser = "regex_parser"
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`(?s)Inappropriate value for attribute "parsers": element 0: attribute.*"regex_parser_options" is required.`),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
						parsers = [
							{
								label = "error logs"
								parser = "regex_parser"
								regex_parser_options = {
								}
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`(?s)Inappropriate value for attribute "parsers": element 0: attribute.*"regex_parser_options": attribute "pattern" is required.`),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
						parsers = [
							{
								label = "error logs"
								parser = "regex_parser"
								regex_parser_options = {
									pattern = ""
								}
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`(?s)Attribute parsers\[0\].regex_parser_options.pattern string length must be at.*least 1, got: 0`),
			},

			// timestamp_parser
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
						parsers = [
							{
								label = "error logs"
								parser = "timestamp_parser"
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`(?s)Inappropriate value for attribute "parsers": element 0: attribute.*"timestamp_parser_options" is required.`),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
						parsers = [
							{
								label = "error logs"
								parser = "timestamp_parser"
								timestamp_parser_options = {
								}
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`(?s)Inappropriate value for attribute "parsers": element 0: attribute.*"timestamp_parser_options": attribute "format" is required.`),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
						parsers = [
							{
								label = "error logs"
								parser = "timestamp_parser"
								timestamp_parser_options = {
									format = ""
								}
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`(?s)Attribute parsers\[0\].timestamp_parser_options.format string length must be at.*least 1, got: 0`),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
						parsers = [
							{
								label = "error logs"
								parser = "timestamp_parser"
								timestamp_parser_options = {
									format = "Custom"
									custom_format = ""
								}
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`(?s)Attribute parsers\[0\].timestamp_parser_options.custom_format string length.*must be at least 1, got: 0`),
			},

			// parser with empty target
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "with_empty_target" {
						title = "custom regex parser title"
						description = "custom regex parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						target_field = ""
						parsers = [
							{
								parser = "regex_parser"
								regex_parser_options = {
									pattern = "\\d{3}-\\d{2}-\\d{3}"
								}
							}
						]
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_sequentially_processor.with_empty_target", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_sequentially_processor.with_empty_target", map[string]any{
						"pipeline_id":                            "#mezmo_pipeline.test_parent.id",
						"title":                                  "custom regex parser title",
						"description":                            "custom regex parser desc",
						"generation_id":                          "0",
						"inputs.#":                               "0",
						"field":                                  ".something",
						"target_field":                           "",
						"parsers.#":                              "1",
						"parsers.0.parser":                       "regex_parser",
						"parsers.0.regex_parser_options.pattern": "\\d{3}-\\d{2}-\\d{3}",
						"parsers.0.regex_parser_options.case_sensitive": "true",
						"parsers.0.regex_parser_options.multiline":      "false",
						"parsers.0.regex_parser_options.match_newline":  "false",
						"parsers.0.regex_parser_options.crlf_newline":   "false",
						"parsers.0.output_name":                         regexp.MustCompile(".+"),
					}),
				),
			},
			// with explicit target
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "with_target" {
						title = "custom regex parser title"
						description = "custom regex parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						target_field = ".data_parsed"
						parsers = [
							{
								parser = "regex_parser"
								regex_parser_options = {
									pattern = "\\d{3}-\\d{2}-\\d{3}"
								}
							}
						]
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_sequentially_processor.with_target", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_sequentially_processor.with_target", map[string]any{
						"pipeline_id":                            "#mezmo_pipeline.test_parent.id",
						"title":                                  "custom regex parser title",
						"description":                            "custom regex parser desc",
						"generation_id":                          "0",
						"inputs.#":                               "0",
						"field":                                  ".something",
						"target_field":                           ".data_parsed",
						"parsers.#":                              "1",
						"parsers.0.parser":                       "regex_parser",
						"parsers.0.regex_parser_options.pattern": "\\d{3}-\\d{2}-\\d{3}",
						"parsers.0.regex_parser_options.case_sensitive": "true",
						"parsers.0.regex_parser_options.multiline":      "false",
						"parsers.0.regex_parser_options.match_newline":  "false",
						"parsers.0.regex_parser_options.crlf_newline":   "false",
						"parsers.0.output_name":                         regexp.MustCompile(".+"),
					}),
				),
			},
			// Create regex parser - default options for regex
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "with_regex" {
						title = "custom regex parser title"
						description = "custom regex parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parsers = [
							{
								parser = "regex_parser"
								regex_parser_options = {
									pattern = "\\d{3}-\\d{2}-\\d{3}"
								}
							}
						]
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_sequentially_processor.with_regex", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_sequentially_processor.with_regex", map[string]any{
						"pipeline_id":                            "#mezmo_pipeline.test_parent.id",
						"title":                                  "custom regex parser title",
						"description":                            "custom regex parser desc",
						"generation_id":                          "0",
						"inputs.#":                               "0",
						"field":                                  ".something",
						"target_field":                           "",
						"parsers.#":                              "1",
						"parsers.0.parser":                       "regex_parser",
						"parsers.0.regex_parser_options.pattern": "\\d{3}-\\d{2}-\\d{3}",
						"parsers.0.regex_parser_options.case_sensitive": "true",
						"parsers.0.regex_parser_options.multiline":      "false",
						"parsers.0.regex_parser_options.match_newline":  "false",
						"parsers.0.regex_parser_options.crlf_newline":   "false",
						"parsers.0.output_name":                         regexp.MustCompile(`^.+\..+$`),
						"unmatched":                                     regexp.MustCompile(`^.+\..+$`),
					}),
				),
			},
			// Create regex parser - custom options for regex
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "with_regex2" {
						title = "custom regex parser title"
						description = "custom regex parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parsers = [
							{
								parser = "regex_parser"
								regex_parser_options = {
									pattern = "\\d{3}-\\d{2}-\\d{3}"
									multiline = true
									case_sensitive = false
									match_newline = true
									crlf_newline = true
								}
							}
						]
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_sequentially_processor.with_regex2", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_sequentially_processor.with_regex2", map[string]any{
						"pipeline_id":                            "#mezmo_pipeline.test_parent.id",
						"title":                                  "custom regex parser title",
						"description":                            "custom regex parser desc",
						"generation_id":                          "0",
						"inputs.#":                               "0",
						"field":                                  ".something",
						"target_field":                           "",
						"parsers.#":                              "1",
						"parsers.0.parser":                       "regex_parser",
						"parsers.0.regex_parser_options.pattern": "\\d{3}-\\d{2}-\\d{3}",
						"parsers.0.regex_parser_options.case_sensitive": "false",
						"parsers.0.regex_parser_options.multiline":      "true",
						"parsers.0.regex_parser_options.match_newline":  "true",
						"parsers.0.regex_parser_options.crlf_newline":   "true",
						"parsers.0.output_name":                         regexp.MustCompile(`^.+\..+$`),
						"unmatched":                                     regexp.MustCompile(`^.+\..+$`),
					}),
				),
			},
			// Create CSV parser with field names
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "csv_parser" {
						title = "custom csv parser title"
						description = "custom csv parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parsers = [
							{
								parser = "csv_row"
								csv_row_options = {
									field_names = ["field1", "field2"]
								}
							}
						]
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_sequentially_processor.csv_parser", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_sequentially_processor.csv_parser", map[string]any{
						"pipeline_id":      "#mezmo_pipeline.test_parent.id",
						"title":            "custom csv parser title",
						"description":      "custom csv parser desc",
						"generation_id":    "0",
						"inputs.#":         "0",
						"field":            ".something",
						"target_field":     "",
						"parsers.#":        "1",
						"parsers.0.parser": "csv_row",
						"parsers.0.csv_row_options.field_names.#": "2",
						"parsers.0.csv_row_options.field_names.0": "field1",
						"parsers.0.csv_row_options.field_names.1": "field2",
						"parsers.0.output_name":                   regexp.MustCompile(`^.+\..+$`),
						"unmatched":                               regexp.MustCompile(`^.+\..+$`),
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "import_target" {
						title = "custom csv parser title"
						description = "custom csv parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parsers = [
							{
								parser = "csv_row"
								csv_row_options = {
									field_names = ["field1", "field2"]
								}
							}
						]
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_parse_sequentially_processor.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_parse_sequentially_processor.csv_parser"),
				ImportStateVerify: true,
			},

			// Create multiple parsers
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "multiple_parsers" {
						title = "create parsers title"
						description = "create parsers desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".app"
						target_field = ".app_parsed"
						parsers = [
							{
								parser = "apache_log"
								apache_log_options = {
									format = "combined"
								}
							},
							{
								parser = "nginx_log"
								nginx_log_options = {
									format = "error"
								}
							},
							{
								parser = "cef_log"
							}
						]
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_sequentially_processor.multiple_parsers", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_sequentially_processor.multiple_parsers", map[string]any{
						"pipeline_id":                         "#mezmo_pipeline.test_parent.id",
						"title":                               "create parsers title",
						"description":                         "create parsers desc",
						"generation_id":                       "0",
						"inputs.#":                            "0",
						"field":                               ".app",
						"target_field":                        ".app_parsed",
						"parsers.#":                           "3",
						"parsers.0.parser":                    "apache_log",
						"parsers.0.apache_log_options.format": "combined",
						"parsers.0.apache_log_options.timestamp_format": "%d/%b/%Y:%T %z",
						"parsers.0.output_name":                         regexp.MustCompile(`^.+\..+$`),
						"parsers.1.parser":                              "nginx_log",
						"parsers.1.nginx_log_options.format":            "error",
						"parsers.1.nginx_log_options.timestamp_format":  "%Y/%m/%d %H:%M:%S",
						"parsers.1.output_name":                         regexp.MustCompile(`^.+\..+$`),
						"parsers.2.parser":                              "cef_log",
						"parsers.2.cef_log_options":                     nil,
						"parsers.2.output_name":                         regexp.MustCompile(`^.+\..+$`),
						"unmatched":                                     regexp.MustCompile(`^.+\..+$`),
					}),
				),
			},
			// Update parser
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "multiple_parsers" {
						title = "create parsers title"
						description = "create parsers desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						field = ".app"
						target_field = ".app_parsed"
						parsers = [
							{
								parser = "regex_parser"
								regex_parser_options = {
									pattern = "(?P<number>[0-9]*)(?P<word>\\w*)"
								}
							},
							{
								parser = "nginx_log"
								nginx_log_options = {
									format = "error"
								}
							},
							{
								parser = "cef_log"
							}
						]
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_sequentially_processor.multiple_parsers", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_sequentially_processor.multiple_parsers", map[string]any{
						"pipeline_id":                            "#mezmo_pipeline.test_parent.id",
						"title":                                  "create parsers title",
						"description":                            "create parsers desc",
						"generation_id":                          "1",
						"inputs.#":                               "1",
						"inputs.0":                               "#mezmo_http_source.my_source.id",
						"field":                                  ".app",
						"target_field":                           ".app_parsed",
						"parsers.#":                              "3",
						"parsers.0.parser":                       "regex_parser",
						"parsers.0.regex_parser_options.pattern": "(?P<number>[0-9]*)(?P<word>\\w*)",
						"parsers.0.output_name":                  regexp.MustCompile(`^.+\..+$`),
						"parsers.1.parser":                       "nginx_log",
						"parsers.1.nginx_log_options.format":     "error",
						"parsers.1.nginx_log_options.timestamp_format": "%Y/%m/%d %H:%M:%S",
						"parsers.1.output_name":                        regexp.MustCompile(`^.+\..+$`),
						"parsers.2.parser":                             "cef_log",
						"parsers.2.cef_log_options":                    nil,
						"parsers.2.output_name":                        regexp.MustCompile(`^.+\..+$`),
						"unmatched":                                    regexp.MustCompile(`^.+\..+$`),
					}),
				),
			},
			// Server side validation on create
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "json_parser" {
						title = "parser title"
						description = "parser desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parsers = [
							{
								parser = "json_parser"
							},
							{
								parser = "json_parser"
							}
						]
					}`,
				ExpectError: regexp.MustCompile("(?s)parser.*cannot be defined more than once with the same.*options."),
			},
			// Error: server-side validation on update
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_processor" "multiple_parsers" {
						title = "create parsers title"
						description = "create parsers desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						field = ".app"
						target_field = ".app_parsed"
						parsers = [
							{
								parser = "regex_parser"
								regex_parser_options = {
									pattern = "(?P<number>[0-9]*)(?P<word>\\w*)"
								}
							},
							{
								parser = "regex_parser"
								regex_parser_options = {
									pattern = "(?P<number>[0-9]*)(?P<word>\\w*)"
								}
							},
							{
								parser = "nginx_log"
								nginx_log_options = {
									format = "error"
								}
							}
						]
					}`,
				ExpectError: regexp.MustCompile("(?s)parser.*cannot be defined more than once with the same.*options."),
			},
			// confirm manually deleted resources are recreated
			{
				Config: GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_parse_sequentially_processor" "test_processor" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					inputs 			= []
					field 			= ".something"
					parsers 		= [
						{
							parser 					= "csv_row"
							csv_row_options = {
								field_names 	= ["field1", "field2"]
							}
						}
					]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_sequentially_processor.test_processor", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_parse_sequentially_processor.test_processor", "title", "new title"),
					// delete the resource
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_parse_sequentially_processor.test_processor",
					),
				),
				// verify resource will be re-created after refresh
				ExpectNonEmptyPlan: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
