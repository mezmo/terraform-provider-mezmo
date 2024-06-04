package alerts

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestAbsenceAlert_success(t *testing.T) {
	const cacheKey = "absence_alert"
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

			// CREATE Absence Alert for metric, minimal configuration
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_absence_alert" "default_metric" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "source"
						component_id = mezmo_http_source.my_source.id
						inputs = [mezmo_http_source.my_source.id]
						name = "my absence alert"
						event_type = "metric"
						subject = "Absence Alert"
						body = "You received a absence alert"
						ingestion_key = "abc123"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_absence_alert.default_metric", "id", regexp.MustCompile(`[\w-]{36}`),
					),
					StateHasExpectedValues("mezmo_absence_alert.default_metric", map[string]any{
						"pipeline_id":             "#mezmo_pipeline.test_parent.id",
						"inputs.#":                "1",
						"inputs.0":                "#mezmo_http_source.my_source.id",
						"active":                  "true",
						"component_id":            "#mezmo_http_source.my_source.id",
						"component_kind":          "source",
						"event_type":              "metric",
						"ingestion_key":           "abc123",
						"body":                    "You received a absence alert",
						"name":                    "my absence alert",
						"severity":                "INFO",
						"style":                   "static",
						"subject":                 "Absence Alert",
						"window_duration_minutes": "5",
						"window_type":             "tumbling",
					}),
				),
			},

			// UPDATE Absence Alert
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_absence_alert" "default_metric" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "source"
						component_id = mezmo_http_source.my_source.id
						inputs = [mezmo_http_source.my_source.id]
						name = "updated name"
						description = "updated description"
						event_type = "metric"
						window_type = "sliding"
						window_duration_minutes = 10
						group_by = [".other"]
						severity = "WARNING"
						subject = "updated subject"
						body = "updated body"
						ingestion_key = "abc123"
						active = false
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_absence_alert.default_metric", map[string]any{
						"active":                  "false",
						"description":             "updated description",
						"group_by.#":              "1",
						"group_by.0":              ".other",
						"body":                    "updated body",
						"name":                    "updated name",
						"severity":                "WARNING",
						"subject":                 "updated subject",
						"window_duration_minutes": "10",
						"window_type":             "sliding",
					}),
				),
			},

			// CREATE Absence Alert for log, minimal configuration
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_absence_alert" "default_log" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "source"
						component_id = mezmo_http_source.my_source.id
						inputs = [mezmo_http_source.my_source.id]
						name = "my absence alert"
						event_type = "log"
						event_timestamp = ".timestamp"
						group_by = [".name", ".namespace", ".tags"]
						subject = "Absence Alert for Log event"
						body = "You received a absence alert for a Log event"
						ingestion_key = "abc123"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_absence_alert.default_log", "id", regexp.MustCompile(`[\w-]{36}`),
					),
					StateHasExpectedValues("mezmo_absence_alert.default_log", map[string]any{
						"pipeline_id":             "#mezmo_pipeline.test_parent.id",
						"inputs.#":                "1",
						"inputs.0":                "#mezmo_http_source.my_source.id",
						"active":                  "true",
						"component_id":            "#mezmo_http_source.my_source.id",
						"component_kind":          "source",
						"event_type":              "log",
						"event_timestamp":         ".timestamp",
						"group_by.#":              "3",
						"group_by.0":              ".name",
						"group_by.1":              ".namespace",
						"group_by.2":              ".tags",
						"ingestion_key":           "abc123",
						"name":                    "my absence alert",
						"severity":                "INFO",
						"style":                   "static",
						"subject":                 "Absence Alert for Log event",
						"body":                    "You received a absence alert for a Log event",
						"window_duration_minutes": "5",
						"window_type":             "tumbling",
					}),
				),
			},
		},
	})
}

func TestAbsenceAlert_schema_validation_errors(t *testing.T) {
	const cacheKey = "absence_alert_schema_validation_errors"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// conditionals are not allowed for absence alerts
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_http_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`) + `
					resource "mezmo_absence_alert" "bad_conditional" {
						pipeline_id = mezmo_pipeline.test_parent.id
						component_kind = "source"
						component_id = mezmo_http_source.my_source.id
						inputs = [mezmo_http_source.my_source.id]
						name = "my bad absence alert"
						event_type = "metric"
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
						subject = "Absence Alert for Log event"
						body = "This will not work"
						ingestion_key = "abc123"
					}`,
				ExpectError: regexp.MustCompile(`(?s).*An argument named "conditional" is not expected her`),
			},
		},
	})
}

// For other custom error tests, see `threshold_test.go`. The schema is shared between
// all the alert types, so only 1 file needs to test certain things.
