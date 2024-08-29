package destinations

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/providertest"
)

func TestAccBlackholeDestinationResource(t *testing.T) {
	const cacheKey = "blackhole_destination_resource"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Required field pipeline id
			{
				Config: GetProviderConfig() + `
					resource "mezmo_blackhole_destination" "my_destination" {}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required, but no definition was found"),
			},
			// Create and Read testing
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}`) + `
					resource "mezmo_blackhole_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my destination title"
						description = "my destination description"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_blackhole_destination.my_destination", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_blackhole_destination.my_destination", map[string]any{
						"description":   "my destination description",
						"generation_id": "0",
						"title":         "my destination title",
						"ack_enabled":   "true",
						"inputs.#":      "0",
					}),
				),
			},
			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_blackhole_destination" "import_target" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my destination title"
						description = "my destination description"
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_blackhole_destination.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_blackhole_destination.my_destination"),
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_demo_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						format = "json"
					}
					resource "mezmo_blackhole_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "new title"
						description = "new description"
						ack_enabled = false
						inputs = [mezmo_demo_source.my_source.id]
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_blackhole_destination.my_destination", map[string]any{
						"description":   "new description",
						"generation_id": "1",
						"title":         "new title",
						"ack_enabled":   "false",
						"inputs.#":      "1",
					}),
				),
			},
			// confirm manually deleted resources are recreated
			{
				Config: GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_blackhole_destination" "destination_2" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_blackhole_destination.destination_2", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_blackhole_destination.destination_2", "title", "new title"),
					// delete the resource
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_blackhole_destination.destination_2",
					),
				),
				// verify resource will be re-created after refresh
				ExpectNonEmptyPlan: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
