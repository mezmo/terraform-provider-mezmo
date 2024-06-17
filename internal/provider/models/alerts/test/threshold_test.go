package alerts

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestAccThresholdAlert_success(t *testing.T) {
	const cacheKey = "threshold_alert"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Cache base resources of pipeline and source
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_http_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`),
			},

			// CREATE Threshold Alert for metric, minimal configuration
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_threshold_alert" "default_metric" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "source"
						component_id = mezmo_http_source.my_source.id
						inputs = [mezmo_http_source.my_source.id]
						name = "my threshold alert"
						event_type = "metric"
						operation = "sum"
						conditional = {
							expressions = [
								{
									field = ".event_count"
									operator = "greater"
									value_number = 5000
								}
							],
						}
						subject = "Threshold Alert"
						body = "You received a threshold alert"
						ingestion_key = "abc123"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_threshold_alert.default_metric", "id", regexp.MustCompile(`[\w-]{36}`),
					),
					StateHasExpectedValues("mezmo_threshold_alert.default_metric", map[string]any{
						"pipeline_id":                            "#mezmo_pipeline.test_parent.id",
						"inputs.#":                               "1",
						"inputs.0":                               "#mezmo_http_source.my_source.id",
						"active":                                 "true",
						"component_id":                           "#mezmo_http_source.my_source.id",
						"component_kind":                         "source",
						"conditional.%":                          "3",
						"conditional.expressions.#":              "1",
						"conditional.expressions.0.field":        ".event_count",
						"conditional.expressions.0.operator":     "greater",
						"conditional.expressions.0.value_number": "5000",
						"conditional.logical_operation":          "AND",
						"event_type":                             "metric",
						"ingestion_key":                          "abc123",
						"body":                                   "You received a threshold alert",
						"name":                                   "my threshold alert",
						"operation":                              "sum",
						"severity":                               "INFO",
						"style":                                  "static",
						"subject":                                "Threshold Alert",
						"window_duration_minutes":                "5",
						"window_type":                            "tumbling",
					}),
				),
			},

			// UPDATE Threshold Alert
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_threshold_alert" "default_metric" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "source"
						component_id = mezmo_http_source.my_source.id
						inputs = [mezmo_http_source.my_source.id]
						name = "updated name"
						description = "updated description"
						event_type = "metric"
						operation = "custom"
						script = "function myFunc(a, e, m) { return a }"
						window_type = "sliding"
						window_duration_minutes = 10
						conditional = {
							expressions = [
								{
									field = ".event_count"
									operator = "greater"
									value_number = 5000
								}
							],
						}
						severity = "WARNING"
						subject = "updated subject"
						body = "updated body"
						ingestion_key = "abc123"
						active = false
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_threshold_alert.default_metric", map[string]any{
						"active":                  "false",
						"description":             "updated description",
						"body":                    "updated body",
						"name":                    "updated name",
						"operation":               "custom",
						"script":                  "function myFunc(a, e, m) { return a }",
						"severity":                "WARNING",
						"subject":                 "updated subject",
						"window_duration_minutes": "10",
						"window_type":             "sliding",
					}),
				),
			},

			// CREATE Threshold Alert for log, minimal configuration
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_threshold_alert" "default_log" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "source"
						component_id = mezmo_http_source.my_source.id
						inputs = [mezmo_http_source.my_source.id]
						name = "my threshold alert"
						event_type = "log"
						operation = "custom"
						script = "function myFunc(a, e, m) { return a }"
						event_timestamp = ".timestamp"
						group_by = [".name", ".namespace", ".tags"]
						conditional = {
							expressions = [
								{
									field = ".event_count"
									operator = "greater"
									value_number = 5000
								}
							],
						}
						subject = "Threshold Alert for Log event"
						body = "You received a threshold alert for a Log event"
						ingestion_key = "abc123"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_threshold_alert.default_log", "id", regexp.MustCompile(`[\w-]{36}`),
					),
					StateHasExpectedValues("mezmo_threshold_alert.default_log", map[string]any{
						"pipeline_id":                            "#mezmo_pipeline.test_parent.id",
						"inputs.#":                               "1",
						"inputs.0":                               "#mezmo_http_source.my_source.id",
						"active":                                 "true",
						"component_id":                           "#mezmo_http_source.my_source.id",
						"component_kind":                         "source",
						"conditional.%":                          "3",
						"conditional.expressions.#":              "1",
						"conditional.expressions.0.field":        ".event_count",
						"conditional.expressions.0.operator":     "greater",
						"conditional.expressions.0.value_number": "5000",
						"conditional.logical_operation":          "AND",
						"event_type":                             "log",
						"event_timestamp":                        ".timestamp",
						"group_by.#":                             "3",
						"group_by.0":                             ".name",
						"group_by.1":                             ".namespace",
						"group_by.2":                             ".tags",
						"ingestion_key":                          "abc123",
						"name":                                   "my threshold alert",
						"operation":                              "custom",
						"script":                                 "function myFunc(a, e, m) { return a }",
						"severity":                               "INFO",
						"style":                                  "static",
						"subject":                                "Threshold Alert for Log event",
						"body":                                   "You received a threshold alert for a Log event",
						"window_duration_minutes":                "5",
						"window_type":                            "tumbling",
					}),
				),
			},
		},
	})
}

// WARNING: Mixing test that throw normal schema validation errors cannot be mixed with
// tests where we bubble up errors with diag. This tends to cause teardown errors and
// may be a bug.

func TestThresholdAlert_root_required_errors(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		ErrorCheck: CheckMultipleErrors([]string{
			`The argument "pipeline_id" is required`,
			`The argument "component_kind" is required`,
			`The argument "component_id" is required`,
			`The argument "name" is required`,
			`The argument "event_type" is required`,
			`The argument "operation" is required`,
			`The argument "subject" is required`,
			`The argument "body" is required`,
			`The argument "ingestion_key" is required`,
			`The argument "conditional" is required`,
		}),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_http_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}
					resource "mezmo_threshold_alert" "bad_pipeline" {
					}`,
			},
		},
	})
}

func TestThresholdAlert_schema_validation_errors(t *testing.T) {
	const cacheKey = "threshold_alert_schema_validation_errors"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// only Non_Change_Operator_Labels are allowed in conditional
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_http_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`) + `
					resource "mezmo_threshold_alert" "bad_conditional" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "source"
						component_id = mezmo_http_source.my_source.id
						inputs = [mezmo_http_source.my_source.id]
						name = "my threshold alert"
						event_type = "metric"
						operation = "sum"
						conditional = {
							expressions = [
								{
									field = ".event_count"
									operator = "percent_change_greater"
									value_number = 10
								}
							],
						}
						subject = "Threshold Alert for Log event"
						body = "You received a threshold alert for a Log event"
						ingestion_key = "abc123"
					}`,
				ExpectError: regexp.MustCompile("(?s).*operator value must be one of.*"),
			},
			// invalid pipeline_id
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_threshold_alert" "bad_pipeline" {
						pipeline_id = ""
						component_kind = "source"
						component_id = mezmo_http_source.my_source.id
						inputs = [mezmo_http_source.my_source.id]
						name = "my threshold alert"
						event_type = "metric"
						operation = "sum"
						conditional = {
							expressions = [
								{
									field = ".event_count"
									operator = "greater"
									value_number = 10
								}
							],
						}
						subject = "Threshold Alert for Log event"
						body = "You received a threshold alert for a Log event"
						ingestion_key = "abc123"
					}`,
				ExpectError: regexp.MustCompile("(?s).*Attribute pipeline_id string length must be at least 1"),
			},
			// invalid component_id
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_threshold_alert" "bad_component_id" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "source"
						component_id = ""
						inputs = [mezmo_http_source.my_source.id]
						name = "my threshold alert"
						event_type = "metric"
						operation = "sum"
						conditional = {
							expressions = [
								{
									field = ".event_count"
									operator = "greater"
									value_number = 10
								}
							],
						}
						subject = "Threshold Alert for Log event"
						body = "You received a threshold alert for a Log event"
						ingestion_key = "abc123"
					}`,
				ExpectError: regexp.MustCompile("(?s).*Attribute component_id string length must be at least 1"),
			},
			// invalid component_kind
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_threshold_alert" "bad_component_kind" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "NOPE"
						component_id = mezmo_http_source.my_source.id
						inputs = [mezmo_http_source.my_source.id]
						name = "my threshold alert"
						event_type = "metric"
						operation = "sum"
						conditional = {
							expressions = [
								{
									field = ".event_count"
									operator = "greater"
									value_number = 10
								}
							],
						}
						subject = "Threshold Alert for Log event"
						body = "You received a threshold alert for a Log event"
						ingestion_key = "abc123"
					}`,
				ExpectError: regexp.MustCompile(`(?s).*Attribute component_kind value must be one of: \["source" "transform"\]`),
			},
			// Missing inputs
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_threshold_alert" "missing_inputs" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "source"
						component_id = mezmo_http_source.my_source.id
						inputs = []
						name = "my threshold alert"
						event_type = "metric"
						operation = "sum"
						conditional = {
							expressions = [
								{
									field = ".event_count"
									operator = "greater"
									value_number = 10
								}
							],
						}
						subject = "Threshold Alert for Log event"
						body = "You received a threshold alert for a Log event"
						ingestion_key = "abc123"
					}`,
				ExpectError: regexp.MustCompile(`(?s).*Attribute inputs list must contain at least 1 elements`),
			},
			// Inputs have bad length
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_threshold_alert" "blank_inputs" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "source"
						component_id = mezmo_http_source.my_source.id
						inputs = [""]
						name = "my threshold alert"
						event_type = "metric"
						operation = "sum"
						conditional = {
							expressions = [
								{
									field = ".event_count"
									operator = "greater"
									value_number = 10
								}
							],
						}
						subject = "Threshold Alert for Log event"
						body = "You received a threshold alert for a Log event"
						ingestion_key = "abc123"
					}`,
				ExpectError: regexp.MustCompile(`(?s).*Attribute inputs\[0\] string length must be at least 1`),
			},
			// FIXME: Something causes a teardown failure if this is uncommented, even though nothing
			// FIXME: should be stored in state when running these failure tests.
			//
			// // Inputs have duplicates
			// {
			// 	Config: GetCachedConfig(cacheKey) + `
			// 		resource "mezmo_threshold_alert" "dupe_inputs" {
			// 			pipeline_id = mezmo_pipeline.test_parent.id
			// 			component_kind = "source"
			// 			component_id = mezmo_http_source.my_source.id
			// 			inputs = [
			// 				mezmo_http_source.my_source.id,
			// 				mezmo_http_source.my_source.id
			// 			]
			// 			name = "my threshold alert"
			// 			event_type = "metric"
			// 			operation = "sum"
			// 			conditional = {
			// 				expressions = [
			// 					{
			// 						field = ".event_count"
			// 						operator = "greater"
			// 						value_number = 10
			// 					}
			// 				],
			// 			}
			// 			subject = "Threshold Alert for Log event"
			// 			body = "You received a threshold alert for a Log event"
			// 			ingestion_key = "abc123"
			// 		}`,
			// 	ExpectError: regexp.MustCompile(`(?s).*inputs = \[.*This attribute contains duplicate values of`),
			// },
			// name is invalid
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_threshold_alert" "bad_severity" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "source"
						component_id = mezmo_http_source.my_source.id
						inputs = [mezmo_http_source.my_source.id]
						name = ""
						event_type = "metric"
						operation = "sum"
						conditional = {
							expressions = [
								{
									field = ".event_count"
									operator = "greater"
									value_number = 10
								}
							],
						}
						subject = "Threshold Alert for Log event"
						body = "You received a threshold alert for a Log event"
						ingestion_key = "abc123"
					}`,
				ExpectError: regexp.MustCompile("(?s).*Attribute name string length must be at least 1"),
			},
			// bad event_type
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_threshold_alert" "bad_severity" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "source"
						component_id = mezmo_http_source.my_source.id
						inputs = [mezmo_http_source.my_source.id]
						name = "my name"
						event_type = "notlog"
						operation = "sum"
						conditional = {
							expressions = [
								{
									field = ".event_count"
									operator = "greater"
									value_number = 10
								}
							],
						}
						subject = "Threshold Alert for Log event"
						body = "You received a threshold alert for a Log event"
						ingestion_key = "abc123"
					}`,
				ExpectError: regexp.MustCompile(`(?s).*Attribute event_type value must be one of: \["log" "metric"\]`),
			},
			// bad window_type
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_threshold_alert" "bad_severity" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "source"
						component_id = mezmo_http_source.my_source.id
						inputs = [mezmo_http_source.my_source.id]
						name = "my name"
						event_type = "metric"
						operation = "sum"
						conditional = {
							expressions = [
								{
									field = ".event_count"
									operator = "greater"
									value_number = 10
								}
							],
						}
						window_type = "badwindow"
						subject = "Threshold Alert for Log event"
						body = "You received a threshold alert for a Log event"
						ingestion_key = "abc123"
					}`,
				ExpectError: regexp.MustCompile(`(?s).*Attribute window_type value must be one of: \["tumbling" "sliding"\]`),
			},
			// invalid severity
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_threshold_alert" "bad_severity" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "source"
						component_id = mezmo_http_source.my_source.id
						inputs = [mezmo_http_source.my_source.id]
						name = "my threshold alert"
						event_type = "metric"
						operation = "sum"
						conditional = {
							expressions = [
								{
									field = ".event_count"
									operator = "greater"
									value_number = 10
								}
							],
						}
						subject = "Threshold Alert for Log event"
						body = "You received a threshold alert for a Log event"
						ingestion_key = "abc123"
						severity = "invalid"
					}`,
				ExpectError: regexp.MustCompile("(?s).*Attribute severity value must be one of.*"),
			},
			// invalid operation
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_threshold_alert" "bad_operation" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "source"
						component_id = mezmo_http_source.my_source.id
						inputs = [mezmo_http_source.my_source.id]
						name = "my threshold alert"
						event_type = "metric"
						operation = "NOPE"
						conditional = {
							expressions = [
								{
									field = ".event_count"
									operator = "greater"
									value_number = 10
								}
							],
						}
						subject = "Threshold Alert for Log event"
						body = "You received a threshold alert for a Log event"
						ingestion_key = "abc123"
					}`,
				ExpectError: regexp.MustCompile(`(?s).Attribute operation value must be one of:`),
			},
			// invalid style
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_threshold_alert" "bad_style" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "source"
						component_id = mezmo_http_source.my_source.id
						inputs = [mezmo_http_source.my_source.id]
						name = "my threshold alert"
						event_type = "metric"
						operation = "sum"
						conditional = {
							expressions = [
								{
									field = ".event_count"
									operator = "greater"
									value_number = 10
								}
							],
						}
						style = "NOPE"
						subject = "Threshold Alert for Log event"
						body = "You received a threshold alert for a Log event"
						ingestion_key = "abc123"
					}`,
				ExpectError: regexp.MustCompile(`(?s).Attribute style value must be one of:`),
			},
		},
	})
}

func TestThresholdAlert_custom_errors(t *testing.T) {
	const cacheKey = "threshold_custom_errors"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// custom operation requires script
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_http_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`) + `
					resource "mezmo_threshold_alert" "bad_alert" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "source"
						component_id = mezmo_http_source.my_source.id
						inputs = [mezmo_http_source.my_source.id]
						name = "my threshold alert"
						event_type = "metric"
						operation = "custom"
						conditional = {
							expressions = [
								{
									field = ".event_count"
									operator = "greater"
									value_number = 5000
								}
							],
						}
						subject = "Threshold Alert for Log event"
						body = "You received a threshold alert for a Log event"
						ingestion_key = "abc123"
					}`,
				ExpectError: regexp.MustCompile("(?s).*A 'custom' operation requires a valid JS `script` function"),
			},
			// script cannot be provided without a custom operation.
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_threshold_alert" "bad_alert" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "source"
						component_id = mezmo_http_source.my_source.id
						inputs = [mezmo_http_source.my_source.id]
						name = "my threshold alert"
						event_type = "metric"
						operation = "sum"
						script = "function myFunc (a, e, m) { return a }"
						conditional = {
							expressions = [
								{
									field = ".event_count"
									operator = "greater"
									value_number = 5000
								}
							],
						}
						subject = "Threshold Alert for Log event"
						body = "You received a threshold alert for a Log event"
						ingestion_key = "abc123"
					}`,
				ExpectError: regexp.MustCompile("(?s).*`script` cannot be set when `operation` is not 'custom'"),
			},
			// log type requires a custom operation.
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_threshold_alert" "bad_alert" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "source"
						component_id = mezmo_http_source.my_source.id
						inputs = [mezmo_http_source.my_source.id]
						name = "my threshold alert"
						event_type = "log"
						operation = "sum"
						conditional = {
							expressions = [
								{
									field = ".event_count"
									operator = "greater"
									value_number = 5000
								}
							],
						}
						subject = "Threshold Alert for Log event"
						body = "You received a threshold alert for a Log event"
						ingestion_key = "abc123"
					}`,
				ExpectError: regexp.MustCompile("(?s).*A 'log' event type requires a 'custom' `operation` and a valid JS `script`"),
			},
			// metric type disallows event_timestamp
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_threshold_alert" "bad_alert" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "source"
						component_id = mezmo_http_source.my_source.id
						inputs = [mezmo_http_source.my_source.id]
						name = "my threshold alert"
						event_type = "metric"
						operation = "sum"
						event_timestamp = ".timestamp"
						conditional = {
							expressions = [
								{
									field = ".event_count"
									operator = "greater"
									value_number = 5000
								}
							],
						}
						subject = "Threshold Alert for Log event"
						body = "You received a threshold alert for a Log event"
						ingestion_key = "abc123"
					}`,
				ExpectError: regexp.MustCompile("(?s).*A 'metric' event type cannot have an `event_timestamp` field"),
			},
		},
	})
}
