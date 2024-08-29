package alerts

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/providertest"
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
						alert_payload = {
							service = {
								name = "log_analysis"
								subject = "Threshold Alert"
								body = "You received a threshold alert"
								ingestion_key = "abc123"
								severity = "INFO"
							}
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_threshold_alert.default_metric", "id", regexp.MustCompile(`[\w-]{36}`),
					),
					resource.TestCheckNoResourceAttr("mezmo_threshold_alert.default_metric", "alert_payload.service.auth"),
					resource.TestCheckNoResourceAttr("mezmo_threshold_alert.default_metric", "alert_payload.service.headers"),
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
						"name":                                   "my threshold alert",
						"operation":                              "sum",
						"window_duration_minutes":                "5",
						"window_type":                            "tumbling",
						"alert_payload.service.name":             "log_analysis",
						"alert_payload.service.ingestion_key":    "abc123",
						"alert_payload.service.body":             "You received a threshold alert",
						"alert_payload.service.severity":         "INFO",
						"alert_payload.service.subject":          "Threshold Alert",
						"alert_payload.throttling.window_secs":   "60",
						"alert_payload.throttling.threshold":     "1",
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
						alert_payload = {
							service = {
								name = "log_analysis"
								severity = "WARNING"
								subject = "updated subject"
								body = "updated body"
								ingestion_key = "abc123"
							}
							throttling = {
                window_secs = 300
								threshold = 5
              }
						}
						active = false
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("mezmo_threshold_alert.default_metric", "alert_payload.service.auth"),
					resource.TestCheckNoResourceAttr("mezmo_threshold_alert.default_metric", "alert_payload.service.headers"),
					StateHasExpectedValues("mezmo_threshold_alert.default_metric", map[string]any{
						"active":                               "false",
						"description":                          "updated description",
						"name":                                 "updated name",
						"operation":                            "custom",
						"script":                               "function myFunc(a, e, m) { return a }",
						"window_duration_minutes":              "10",
						"window_type":                          "sliding",
						"alert_payload.service.name":           "log_analysis",
						"alert_payload.service.ingestion_key":  "abc123",
						"alert_payload.service.body":           "updated body",
						"alert_payload.service.severity":       "WARNING",
						"alert_payload.service.subject":        "updated subject",
						"alert_payload.throttling.window_secs": "300",
						"alert_payload.throttling.threshold":   "5",
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
						alert_payload = {
							service = {
								name = "slack"
								uri = "http://google.com/our_slack_api"
								message_text = "Alert: Log event count has exceeded threshold"
							}
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_threshold_alert.default_log", "id", regexp.MustCompile(`[\w-]{36}`),
					),
					resource.TestCheckNoResourceAttr("mezmo_threshold_alert.default_log", "alert_payload.service.auth"),
					resource.TestCheckNoResourceAttr("mezmo_threshold_alert.default_log", "alert_payload.service.headers"),
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
						"name":                                   "my threshold alert",
						"operation":                              "custom",
						"script":                                 "function myFunc(a, e, m) { return a }",
						"window_duration_minutes":                "5",
						"window_type":                            "tumbling",
						"alert_payload.service.name":             "slack",
						"alert_payload.service.message_text":     "Alert: Log event count has exceeded threshold",
						"alert_payload.service.uri":              "http://google.com/our_slack_api",
					}),
				),
			},
			// CREATE alert with webhook auth payload
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_threshold_alert" "webhook_auth_alert" {
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
						alert_payload = {
							service = {
								name = "webhook"
								uri = "http://example.com/our_webhook_api"
								message_text = "Alert: Log event count has exceeded threshold"
								auth = {
									strategy = "basic"
									user = "my_user"
									password = "my_password"
								}
								headers = {
                  "x-custom-header" = "header_value"
                }
							}
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_threshold_alert.webhook_auth_alert", map[string]any{
						"alert_payload.service.auth.password":           "my_password",
						"alert_payload.service.auth.strategy":           "basic",
						"alert_payload.service.auth.user":               "my_user",
						"alert_payload.service.headers.%":               "1",
						"alert_payload.service.headers.x-custom-header": "header_value",
					}),
					resource.TestCheckNoResourceAttr("mezmo_threshold_alert.webhook_auth_alert", "alert_payload.service.auth.token"),
				),
			},
			// Update webhook alert and remove auth and headers
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_threshold_alert" "webhook_auth_alert" {
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
						alert_payload = {
							service = {
								name = "webhook"
								uri = "http://example.com/our_webhook_api"
								message_text = "Alert: Log event count has exceeded threshold"
							}
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("mezmo_threshold_alert.webhook_auth_alert", "alert_payload.service.auth"),
					resource.TestCheckNoResourceAttr("mezmo_threshold_alert.webhook_auth_alert", "alert_payload.service.headers"),
				),
			},
		},
	})
}

// WARNING: Mixing test that throw normal schema validation errors cannot be mixed with
// tests where we bubble up errors with diag. This tends to cause teardown errors and
// may be a bug.

func TestAccThresholdAlert_root_required_errors(t *testing.T) {
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
			`The argument "conditional" is required`,
			`The argument "alert_payload" is required`,
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

func TestAccThresholdAlert_schema_validation_errors(t *testing.T) {
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
						alert_payload = {
							service = {
                name = "slack"
                uri = "http://google.com/our_slack_api"
                message_text = "Alert: Log event count has exceeded threshold"
              }
						}
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
						alert_payload = {
							service = {
                name = "slack"
                uri = "http://google.com/our_slack_api"
                message_text = "Alert: Log event count has exceeded threshold"
              }
						}
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
						alert_payload = {
							service = {
                name = "slack"
                uri = "http://google.com/our_slack_api"
                message_text = "Alert: Log event count has exceeded threshold"
              }
						}
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
						alert_payload = {
							service = {
                name = "slack"
                uri = "http://google.com/our_slack_api"
                message_text = "Alert: Log event count has exceeded threshold"
              }
						}
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
						alert_payload = {
							service = {
                name = "slack"
                uri = "http://google.com/our_slack_api"
                message_text = "Alert: Log event count has exceeded threshold"
              }
						}
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
						alert_payload = {
							service = {
                name = "slack"
                uri = "http://google.com/our_slack_api"
                message_text = "Alert: Log event count has exceeded threshold"
              }
						}
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
						alert_payload = {
							service = {
                name = "slack"
                uri = "http://google.com/our_slack_api"
                message_text = "Alert: Log event count has exceeded threshold"
              }
						}
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
						alert_payload = {
							service = {
                name = "slack"
                uri = "http://google.com/our_slack_api"
                message_text = "Alert: Log event count has exceeded threshold"
              }
						}
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
						alert_payload = {
							service = {
                name = "slack"
                uri = "http://google.com/our_slack_api"
                message_text = "Alert: Log event count has exceeded threshold"
              }
						}
					}`,
				ExpectError: regexp.MustCompile(`(?s).*Attribute window_type value must be one of: \["tumbling" "sliding"\]`),
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
						alert_payload = {
							service = {
                name = "slack"
                uri = "http://google.com/our_slack_api"
                message_text = "Alert: Log event count has exceeded threshold"
              }
						}
					}`,
				ExpectError: regexp.MustCompile(`(?s).*Attribute operation value must be one of:`),
			},
			// invalid auth for basic
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_threshold_alert" "bad_operation" {
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
						alert_payload = {
							service = {
                name = "webhook"
                uri = "http://example.com/our_webhook_api"
                message_text = "Alert: Log event count has exceeded threshold"
								auth = {
									strategy = "basic"
								}
              }
						}
					}`,
				ExpectError: regexp.MustCompile(`(?s).*Basic auth requires user and password fields to be defined`),
			},
			// invalid auth for bearer
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_threshold_alert" "bad_operation" {
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
						alert_payload = {
							service = {
                name = "webhook"
                uri = "http://example.com/our_webhook_api"
                message_text = "Alert: Log event count has exceeded threshold"
								auth = {
									strategy = "bearer"
								}
              }
						}
					}`,
				ExpectError: regexp.MustCompile(`(?s).*Bearer auth requires token field to be defined`),
			},
		},
	})
}

func TestAccThresholdAlert_custom_errors(t *testing.T) {
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
						alert_payload = {
							service = {
                name = "slack"
                uri = "http://google.com/our_slack_api"
                message_text = "Alert: Log event count has exceeded threshold"
              }
						}
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
						alert_payload = {
							service = {
                name = "slack"
                uri = "http://google.com/our_slack_api"
                message_text = "Alert: Log event count has exceeded threshold"
              }
						}
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
						alert_payload = {
							service = {
                name = "slack"
                uri = "http://google.com/our_slack_api"
                message_text = "Alert: Log event count has exceeded threshold"
              }
						}
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
						alert_payload = {
							service = {
                name = "slack"
                uri = "http://google.com/our_slack_api"
                message_text = "Alert: Log event count has exceeded threshold"
              }
						}
					}`,
				ExpectError: regexp.MustCompile("(?s).*A 'metric' event type cannot have an `event_timestamp` field"),
			},
		},
	})
}

func TestAccThresholdAlert_slack_payload_errors(t *testing.T) {
	const cacheKey = "slack_payload_errors"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ErrorCheck: CheckMultipleErrors([]string{
			"`uri` is required for the `slack` service",
			"`message_text` is required for the `slack` service",
			"Attribute `summary` is not allowed for service `slack`",
			"Attribute `source` is not allowed for service `slack`",
			"Attribute `routing_key` is not allowed for service `slack`",
			"Attribute `event_action` is not allowed for service `slack`",
			"Attribute `severity` is not allowed for service `slack`",
			"Attribute `subject` is not allowed for service `slack`",
			"Attribute `body` is not allowed for service `slack`",
			"Attribute `ingestion_key` is not allowed for service `slack`",
			"Attribute `auth` is not allowed for service `slack`",
			"Attribute `headers` is not allowed for service `slack`",
			"Attribute `method` is not allowed for service `slack`",
		}),
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
						alert_payload = {
							service = {
                name = "slack"
								summary = "nope"
								source = "nope"
								routing_key = "nope"
								event_action = "nope"
								severity = "INFO"
								subject = "nope"
								body = "nope"
								ingestion_key = "nope"
								auth = {
									strategy = "bearer"
								}
								headers = {}
								method = "post"
              }
						}
					}`,
			},
		},
	})
}

func TestAccThresholdAlert_webhook_payload_errors(t *testing.T) {
	const cacheKey = "webhook_payload_errors"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		ErrorCheck: CheckMultipleErrors([]string{
			"`uri` is required for the `webhook` service",
			"`message_text` is required for the `webhook` service",
			"Attribute `summary` is not allowed for service `webhook`",
			"Attribute `source` is not allowed for service `webhook`",
			"Attribute `routing_key` is not allowed for service `webhook`",
			"Attribute `event_action` is not allowed for service `webhook`",
			"Attribute `severity` is not allowed for service `webhook`",
			"Attribute `subject` is not allowed for service `webhook`",
			"Attribute `body` is not allowed for service `webhook`",
			"Attribute `ingestion_key` is not allowed for service `webhook`",
		}),
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
						alert_payload = {
							service = {
                name = "webhook"
								summary = "nope"
								source = "nope"
								routing_key = "nope"
								event_action = "nope"
								severity = "INFO"
								subject = "nope"
								body = "nope"
								ingestion_key = "nope"
              }
						}
					}`,
			},
		},
	})
}

func TestAccThresholdAlert_pager_duty_payload_errors(t *testing.T) {
	const cacheKey = "pager_duty_payload_errors"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		ErrorCheck: CheckMultipleErrors([]string{
			"`uri` is required for the `pager_duty` service",
			"`summary` is required for the `pager_duty` service",
			"`severity` is required for the `pager_duty` service",
			"`source` is required for the `pager_duty` service",
			"`routing_key` is required for the `pager_duty` service",
			"`event_action` is required for the `pager_duty` service",
			"Attribute `message_text` is not allowed for service `pager_duty`",
			"Attribute `subject` is not allowed for service `pager_duty`",
			"Attribute `body` is not allowed for service `pager_duty`",
			"Attribute `ingestion_key` is not allowed for service `pager_duty`",
			"Attribute `auth` is not allowed for service `pager_duty`",
			"Attribute `headers` is not allowed for service `pager_duty`",
			"Attribute `method` is not allowed for service `pager_duty`",
		}),
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
						alert_payload = {
							service = {
                name = "pager_duty"
								message_text = "nope"
								subject = "nope"
								body = "nope"
								ingestion_key = "nope"
								auth = {
									strategy = "bearer"
								}
								headers = {}
								method = "post"
              }
						}
					}`,
			},
		},
	})
}

func TestAccThresholdAlert_log_analysis_payload_errors(t *testing.T) {
	const cacheKey = "log_analysis_payload_errors"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		ErrorCheck: CheckMultipleErrors([]string{
			"`severity` is required for the `log_analysis` service",
			"`subject` is required for the `log_analysis` service",
			"`body` is required for the `log_analysis` service",
			"`ingestion_key` is required for the `log_analysis` service",
			"Attribute `uri` is not allowed for service `log_analysis`",
			"Attribute `message_text` is not allowed for service `log_analysis`",
			"Attribute `summary` is not allowed for service `log_analysis`",
			"Attribute `source` is not allowed for service `log_analysis`",
			"Attribute `routing_key` is not allowed for service `log_analysis`",
			"Attribute `event_action` is not allowed for service `log_analysis`",
			"Attribute `auth` is not allowed for service `log_analysis`",
			"Attribute `headers` is not allowed for service `log_analysis`",
			"Attribute `method` is not allowed for service `log_analysis`",
		}),
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
						alert_payload = {
							service = {
                name = "log_analysis"
								uri = "http://example.com/webhook"
								message_text = "nope"
								summary = "nope"
								source = "nope"
								routing_key = "nope"
								event_action = "nope"
								auth = {
									strategy = "bearer"
								}
								headers = {}
								method = "post"
              }
						}
					}`,
			},
		},
	})
}
