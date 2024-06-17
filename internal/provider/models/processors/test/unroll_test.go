package processors

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestAccUnrollProcessor(t *testing.T) {
	const cacheKey = "unroll_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Error: pipeline_id is required
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_http_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`) + `
					resource "mezmo_unroll_processor" "my_processor" {
						field = ".nope"
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},

			// Error: `field` is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_unroll_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile("The argument \"field\" is required"),
			},

			// Error: `field` values validates length
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_unroll_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ""
					}`,
				ExpectError: regexp.MustCompile("Attribute field string length must be at least 1"),
			},

			// Create
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_unroll_processor" "my_processor" {
						title = "processor title"
						description = "processor desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".thing1"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_unroll_processor.my_processor", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_unroll_processor.my_processor", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "processor title",
						"description":   "processor desc",
						"generation_id": "0",
						"inputs.#":      "0",
						"field":         ".thing1",
						"values_only":   "true",
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_unroll_processor" "import_target" {
						title = "processor title"
						description = "processor desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".thing1"
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_unroll_processor.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_unroll_processor.my_processor"),
				ImportStateVerify: true,
			},

			// Update
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_unroll_processor" "my_processor" {
						title = "new title"
						description = "new desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						field = ".thing2"
						values_only = false
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_unroll_processor.my_processor", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "new title",
						"description":   "new desc",
						"generation_id": "1",
						"inputs.#":      "1",
						"inputs.0":      "#mezmo_http_source.my_source.id",
						"field":         ".thing2",
						"values_only":   "false",
					}),
				),
			},

			// Error: server-side validation
			{
				Config: GetCachedConfig(cacheKey) + `
				resource "mezmo_unroll_processor" "my_processor" {
					pipeline_id = mezmo_pipeline.test_parent.id
					inputs = []
					field = "not-a-valid-field"
				}`,
				ExpectError: regexp.MustCompile("match pattern"),
			},

			// confirm manually deleted resources are recreated
			{
				Config: GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_unroll_processor" "test_processor" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					inputs 			= []
					field 			= ".thing1"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_unroll_processor.test_processor", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_unroll_processor.test_processor", "title", "new title"),
					// delete the resource
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_unroll_processor.test_processor",
					),
				),
				// verify resource will be re-created after refresh
				ExpectNonEmptyPlan: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
