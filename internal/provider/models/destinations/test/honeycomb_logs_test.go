package destinations

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/providertest"
)

func TestAccHoneycombLogsDestinationResource(t *testing.T) {
	const cacheKey = "honeycomb_destination_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Error: properties are required
			{
				Config: GetProviderConfig() + `
					resource "mezmo_honeycomb_logs_destination" "my_destination" {
						inputs  = ["abc"]
						api_key = "my_key"
						dataset = "hello"
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_honeycomb_logs_destination" "my_destination" {
						inputs      = ["abc"]
						pipeline_id = "pip1"
						dataset     = "hello"
					}`,
				ExpectError: regexp.MustCompile("The argument \"api_key\" is required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_honeycomb_logs_destination" "my_destination" {
						inputs      = ["abc"]
						api_key     = "my_key"
						pipeline_id = "pip1"
					}`,
				ExpectError: regexp.MustCompile("The argument \"dataset\" is required"),
			},

			// Create test defaults
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_http_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`) + `
					resource "mezmo_honeycomb_logs_destination" "my_destination" {
						title = "My destination"
						description = "my destination description"
						pipeline_id = mezmo_pipeline.test_parent.id
						dataset     = "ds1"
						api_key     = "key1"
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_honeycomb_logs_destination.my_destination", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_honeycomb_logs_destination.my_destination", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "My destination",
						"description":   "my destination description",
						"generation_id": "0",
						"ack_enabled":   "true",
						"inputs.#":      "0",
						"dataset":       "ds1",
						"api_key":       "key1",
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_honeycomb_logs_destination" "import_target" {
						title = "My destination"
						description = "my destination description"
						pipeline_id = mezmo_pipeline.test_parent.id
						dataset     = "ds1"
						api_key     = "key1"
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_honeycomb_logs_destination.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_honeycomb_logs_destination.my_destination"),
				ImportStateVerify: true,
			},

			// Update all fields
			{
				Config: GetCachedConfig(cacheKey) + `
				resource "mezmo_honeycomb_logs_destination" "my_destination" {
					title = "my new destination title"
					description = "my new destination description"
					inputs      = [mezmo_http_source.my_source.id]
					pipeline_id = mezmo_pipeline.test_parent.id
					dataset     = "ds2"
					api_key     = "key2"
				}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_honeycomb_logs_destination.my_destination", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "my new destination title",
						"description":   "my new destination description",
						"generation_id": "1",
						"ack_enabled":   "true",
						"inputs.#":      "1",
						"dataset":       "ds2",
						"api_key":       "key2",
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
				resource "mezmo_honeycomb_logs_destination" "test_destination" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					dataset     = "ds3"
					api_key     = "key2"
					inputs 			= [mezmo_http_source.my_source2.id]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_honeycomb_logs_destination.test_destination", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_honeycomb_logs_destination.test_destination", "title", "new title"),
					// verify resource will be re-created after refresh
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_honeycomb_logs_destination.test_destination",
					),
				),
				ExpectNonEmptyPlan: true,
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
