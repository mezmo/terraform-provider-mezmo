package destinations

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestAccDatadogLogsDestinationResource(t *testing.T) {
	cacheKey := "datadog_logs_destination_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Required field site
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}`) + `
					resource "mezmo_datadog_logs_destination" "my_destination" {
						site        = "us3"
						api_key     = "<secret-api-key>"
						compression = "gzip"
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required, but no definition was found"),
			},

			// Required field api_key
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_datadog_logs_destination" "my_destination" {
						site = "us3"
					}`,
				ExpectError: regexp.MustCompile("The argument \"api_key\" is required, but no definition was found"),
			},

			// Required field site
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_datadog_logs_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						api_key     = "<secret-api-key>"
						compression = "gzip"
					}`,
				ExpectError: regexp.MustCompile("The argument \"site\" is required"),
			},

			// Required field compression
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_datadog_logs_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						site        = "us3"
						api_key     = "<secret-api-key>"
					}`,
				ExpectError: regexp.MustCompile("The argument \"compression\" is required"),
			},

			// Site field acceptable values
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_datadog_logs_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						site        = "blah"
						api_key     = "<secret-api-key>"
						compression = "gzip"
					}`,
				ExpectError: regexp.MustCompile("Attribute site value must be one of"),
			},

			// Compression field acceptable values
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_datadog_logs_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						site        = "us3"
						api_key     = "<secret-api-key>"
						compression = "blah"
					}`,
				ExpectError: regexp.MustCompile("Attribute compression value must be one of"),
			},

			// Test defaults with minimal values
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_http_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`) + `
					resource "mezmo_datadog_logs_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						site        = "us3"
						api_key     = "<secret-api-key>"
						compression = "gzip"
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_datadog_logs_destination.my_destination", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_datadog_logs_destination.my_destination", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"generation_id": "0",
						"ack_enabled":   "true",
						"inputs.#":      "0",
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_datadog_logs_destination" "import_target" {
						title       = "my logs destination"
						description = "logs description"
						pipeline_id = mezmo_pipeline.test_parent.id
						site        = "us3"
						api_key     = "<secret-api-key>"
						compression = "gzip"
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_datadog_logs_destination.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_datadog_logs_destination.my_destination"),
				ImportStateVerify: true,
			},

			// Update all fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_datadog_logs_destination" "my_destination" {
						title = "new title"
						description = "new logs description"
						pipeline_id = mezmo_pipeline.test_parent.id
						site        = "us1"
						api_key     = "<new-secret-api-key>"
						compression = "none"
						ack_enabled = false
						inputs = [mezmo_http_source.my_source.id]
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_datadog_logs_destination.my_destination", map[string]any{
						"pipeline_id": "#mezmo_pipeline.test_parent.id",
						"title":       "new title",
						"description": "new logs description",
						"site":        "us1",
						"api_key":     "<new-secret-api-key>",
						"compression": "none",
						"ack_enabled": "false",
						"inputs.#":    "1",
						"inputs.0":    "#mezmo_http_source.my_source.id",
					}),
				),
			},
			// confirm manually deleted resources are recreated
			{
				Config: GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_http_source" "my_source2" {
					pipeline_id = mezmo_pipeline.test_parent2.id
				}
				resource "mezmo_datadog_logs_destination" "test_destination" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					site        = "us3"
					api_key     = "<secret-api-key>"
					compression = "gzip"
					inputs = [mezmo_http_source.my_source2.id]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_datadog_logs_destination.test_destination", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_datadog_logs_destination.test_destination", "title", "new title"),
					// verify resource will be re-created after refresh
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_datadog_logs_destination.test_destination",
					),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
