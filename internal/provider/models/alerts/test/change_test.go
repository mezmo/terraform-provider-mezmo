package alerts

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/providertest"
)

func TestAccChangeAlert_success(t *testing.T) {
	const cacheKey = "change_alert"
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

			// CREATE Change Alert for metric, minimal configuration
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_change_alert" "default_metric" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "source"
						component_id = mezmo_http_source.my_source.id
						inputs = [mezmo_http_source.my_source.id]
						name = "my change alert"
						event_type = "metric"
						operation = "sum"
						conditional = {
							expressions = [
								{
									field = ".some_value"
									operator = "value_change_greater"
									value_number = 500
								}
							],
						}
						alert_payload = {
							service = {
								name = "log_analysis"
								subject = "Change Alert"
								body = "You received a change alert"
								ingestion_key = "abc123"
								severity = "INFO"
							}
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_change_alert.default_metric", "id", regexp.MustCompile(`[\w-]{36}`),
					),
					StateHasExpectedValues("mezmo_change_alert.default_metric", map[string]any{
						"pipeline_id":                            "#mezmo_pipeline.test_parent.id",
						"inputs.#":                               "1",
						"inputs.0":                               "#mezmo_http_source.my_source.id",
						"active":                                 "true",
						"component_id":                           "#mezmo_http_source.my_source.id",
						"component_kind":                         "source",
						"conditional.%":                          "3",
						"conditional.expressions.#":              "1",
						"conditional.expressions.0.field":        ".some_value",
						"conditional.expressions.0.operator":     "value_change_greater",
						"conditional.expressions.0.value_number": "500",
						"conditional.logical_operation":          "AND",
						"event_type":                             "metric",
						"name":                                   "my change alert",
						"operation":                              "sum",
						"window_duration_minutes":                "5",
						"window_type":                            "tumbling",
						"alert_payload.service.name":             "log_analysis",
						"alert_payload.service.ingestion_key":    "abc123",
						"alert_payload.service.body":             "You received a change alert",
						"alert_payload.service.severity":         "INFO",
						"alert_payload.service.subject":          "Change Alert",
						"alert_payload.throttling.window_secs":   "60",
						"alert_payload.throttling.threshold":     "1",
					}),
				),
			},

			// UPDATE Change Alert
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_change_alert" "default_metric" {
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
						group_by = [".other"]
						active = false
						conditional = {
							expressions = [
								{
									field = ".other_value"
									operator = "value_change_less_or_equal"
									value_number = 100
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
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_change_alert.default_metric", map[string]any{
						"active":                                 "false",
						"description":                            "updated description",
						"group_by.#":                             "1",
						"group_by.0":                             ".other",
						"name":                                   "updated name",
						"operation":                              "custom",
						"conditional.expressions.0.field":        ".other_value",
						"conditional.expressions.0.operator":     "value_change_less_or_equal",
						"conditional.expressions.0.value_number": "100",
						"script":                                 "function myFunc(a, e, m) { return a }",
						"window_duration_minutes":                "10",
						"window_type":                            "sliding",
						"alert_payload.service.name":             "log_analysis",
						"alert_payload.service.ingestion_key":    "abc123",
						"alert_payload.service.body":             "updated body",
						"alert_payload.service.severity":         "WARNING",
						"alert_payload.service.subject":          "updated subject",
						"alert_payload.throttling.window_secs":   "300",
						"alert_payload.throttling.threshold":     "5",
					}),
				),
			},

			// CREATE Change Alert for log, minimal configuration
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_change_alert" "default_log" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "source"
						component_id = mezmo_http_source.my_source.id
						inputs = [mezmo_http_source.my_source.id]
						name = "my change alert"
						event_type = "log"
						operation = "custom"
						script = "function myFunc(a, e, m) { return a }"
						event_timestamp = ".timestamp"
						group_by = [".name", ".namespace", ".tags"]
						conditional = {
							expressions = [
								{
									field = ".some_value"
									operator = "percent_change_less"
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_change_alert.default_log", "id", regexp.MustCompile(`[\w-]{36}`),
					),
					StateHasExpectedValues("mezmo_change_alert.default_log", map[string]any{
						"pipeline_id":                            "#mezmo_pipeline.test_parent.id",
						"inputs.#":                               "1",
						"inputs.0":                               "#mezmo_http_source.my_source.id",
						"active":                                 "true",
						"component_id":                           "#mezmo_http_source.my_source.id",
						"component_kind":                         "source",
						"conditional.%":                          "3",
						"conditional.expressions.#":              "1",
						"conditional.expressions.0.field":        ".some_value",
						"conditional.expressions.0.operator":     "percent_change_less",
						"conditional.expressions.0.value_number": "10",
						"conditional.logical_operation":          "AND",
						"event_type":                             "log",
						"event_timestamp":                        ".timestamp",
						"group_by.#":                             "3",
						"group_by.0":                             ".name",
						"group_by.1":                             ".namespace",
						"group_by.2":                             ".tags",
						"name":                                   "my change alert",
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
		},
	})
}

func TestAccChangeAlert_schema_validation_errors(t *testing.T) {
	const cacheKey = "change_alert_schema_validation_errors"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// only Change_Operator_Labels are allowed in conditional
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_http_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`) + `
					resource "mezmo_change_alert" "bad_conditional" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "source"
						component_id = mezmo_http_source.my_source.id
						inputs = [mezmo_http_source.my_source.id]
						name = "my change alert"
						event_type = "metric"
						operation = "sum"
						group_by = [".timestamp"]
						conditional = {
							expressions = [
								{
									field = ".some_value"
									operator = "greater"
									value_number = 500
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
				ExpectError: regexp.MustCompile("(?s).*operator value must be one of.*percent_change_less"),
			},
		},
	})
}

// For other custom error tests, see `threshold_test.go`. The schema is shared between
// all the alert types, so only 1 file needs to test certain things.
