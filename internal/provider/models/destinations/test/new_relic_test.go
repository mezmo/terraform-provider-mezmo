package destinations

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestAccNewRelicDestinationResource(t *testing.T) {
	const cacheKey = "new_relic_destination_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Error: properties are required
			{
				Config: GetProviderConfig() + `
					resource "mezmo_new_relic_destination" "my_destination" {
						inputs      = ["abc"]
						account_id = "acc1"
						license_key = "key1"
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_new_relic_destination" "my_destination" {
						inputs      = ["abc"]
						pipeline_id = "pip1"
						license_key = "key1"
					}`,
				ExpectError: regexp.MustCompile("The argument \"account_id\" is required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_new_relic_destination" "my_destination" {
						inputs      = ["abc"]
						pipeline_id = "pip1"
						account_id = "acc1"
					}`,
				ExpectError: regexp.MustCompile("The argument \"license_key\" is required"),
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
					resource "mezmo_new_relic_destination" "my_destination" {
						title = "My destination"
						description = "my destination description"
						inputs      = [mezmo_http_source.my_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id
						account_id = "acc1"
						license_key = "key1"
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_new_relic_destination.my_destination", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_new_relic_destination.my_destination", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "My destination",
						"description":   "my destination description",
						"generation_id": "0",
						"ack_enabled":   "true",
						"inputs.#":      "1",
						"account_id":    "acc1",
						"license_key":   "key1",
						"api":           "logs",
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_new_relic_destination" "import_target" {
						title = "My destination"
						description = "my destination description"
						inputs      = [mezmo_http_source.my_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id
						account_id = "acc1"
						license_key = "key1"
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_new_relic_destination.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_new_relic_destination.my_destination"),
				ImportStateVerify: true,
			},

			// Update all fields (pass through)
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_new_relic_destination" "my_destination" {
						title = "new title"
						description = "new description"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs      = [mezmo_http_source.my_source.id]
						account_id  = "acc2"
						license_key = "key2"
						api         = "metrics"
						ack_enabled = false
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_new_relic_destination.my_destination", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "new title",
						"description":   "new description",
						"generation_id": "1",
						"inputs.#":      "1",
						"inputs.0":      "#mezmo_http_source.my_source.id",
						"ack_enabled":   "false",
						"account_id":    "acc2",
						"license_key":   "key2",
						"api":           "metrics",
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
				resource "mezmo_new_relic_destination" "test_destination" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					inputs 			= [mezmo_http_source.my_source2.id]
					account_id = "acc1"
					license_key = "key1"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_new_relic_destination.test_destination", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_new_relic_destination.test_destination", "title", "new title"),
					// verify resource will be re-created after refresh
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_new_relic_destination.test_destination",
					),
				),
				ExpectNonEmptyPlan: true,
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
