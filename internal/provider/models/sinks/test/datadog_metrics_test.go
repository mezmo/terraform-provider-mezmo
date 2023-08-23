package sinks

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/providertest"
)

func TestDatadogMetricsSinkResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Required field site
			{
				Config: GetProviderConfig() + `
					resource "mezmo_datadog-metrics_sink" "my_sink" {
						site    = "us3"
						api_key = "<secret-api-key>"
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required, but no definition was found"),
			},

			// Required field api_key
			{
				Config: GetProviderConfig() + `
					resource "mezmo_datadog-metrics_sink" "my_sink" {
						site = "us3"
					}`,
				ExpectError: regexp.MustCompile("The argument \"api_key\" is required, but no definition was found"),
			},

			// Required field site
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_datadog-metrics_sink" "my_sink" {
						pipeline_id = mezmo_pipeline.test_parent.id
						api_key     = "<secret-api-key>"
					}`,
				ExpectError: regexp.MustCompile("The argument \"site\" is required"),
			},

			// Site field acceptable values
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_datadog-metrics_sink" "my_sink" {
						pipeline_id = mezmo_pipeline.test_parent.id
						site        = "blah"
						api_key     = "<secret-api-key>"
					}`,
				ExpectError: regexp.MustCompile("Attribute site value must be one of"),
			},

			// Test defaults
			{
				Config: SetCachedConfig("http_sink_resources", `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_http_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`) + `
					resource "mezmo_datadog-metrics_sink" "my_sink" {
						title       = "my metrics sink"
						description = "metrics description"
						pipeline_id = mezmo_pipeline.test_parent.id
						site        = "us3"
						api_key     = "<secret-api-key>"
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_datadog-metrics_sink.my_sink", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_datadog-metrics_sink.my_sink", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "my metrics sink",
						"description":   "metrics description",
						"generation_id": "0",
						"ack_enabled":   "true",
						"inputs.#":      "0",
					}),
				),
			},

			// Update all fields
			{
				Config: GetCachedConfig("http_sink_resources") + `
					resource "mezmo_datadog-metrics_sink" "my_sink" {
						title = "new title"
						description = "new metrics description"
						pipeline_id = mezmo_pipeline.test_parent.id
						site        = "us1"
						api_key     = "<new-secret-api-key>"
						ack_enabled = false
						inputs = [mezmo_http_source.my_source.id]
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_datadog-metrics_sink.my_sink", map[string]any{
						"pipeline_id": "#mezmo_pipeline.test_parent.id",
						"title":       "new title",
						"description": "new metrics description",
						"site":        "us1",
						"api_key":     "<new-secret-api-key>",
						"ack_enabled": "false",
						"inputs.#":    "1",
						"inputs.0":    "#mezmo_http_source.my_source.id",
					}),
				),
			},
		},
	})
}
