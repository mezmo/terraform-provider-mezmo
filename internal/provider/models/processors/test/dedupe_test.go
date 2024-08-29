package processors

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/providertest"
)

func TestAccDedupeProcessor(t *testing.T) {
	const cacheKey = "dedupe_resources"
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
					resource "mezmo_dedupe_processor" "my_processor" {
						fields = [".nope"]
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},

			// Error: `fields` is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_dedupe_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile("The argument \"fields\" is required"),
			},

			// Error: `fields` array length validation
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_dedupe_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						fields = []
					}`,
				ExpectError: regexp.MustCompile("Attribute fields list must contain at least 1 elements"),
			},

			// Error: `fields` values validates length
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_dedupe_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						fields = [""]
					}`,
				ExpectError: regexp.MustCompile("Attribute fields\\[0\\] string length must be at least 1"),
			},

			// Create with defaults
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_dedupe_processor" "my_processor" {
						title = "dedupe title"
						description = "dedupe desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						fields = [".thing1", ".thing2"]
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_dedupe_processor.my_processor", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_dedupe_processor.my_processor", map[string]any{
						"pipeline_id":      "#mezmo_pipeline.test_parent.id",
						"title":            "dedupe title",
						"description":      "dedupe desc",
						"generation_id":    "0",
						"inputs.#":         "0",
						"number_of_events": "5000",
						"comparison_type":  "Match",
						"fields.#":         "2",
						"fields.0":         ".thing1",
						"fields.1":         ".thing2",
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_dedupe_processor" "import_target" {
						title = "dedupe title"
						description = "dedupe desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						fields = [".thing1", ".thing2"]
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_dedupe_processor.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_dedupe_processor.my_processor"),
				ImportStateVerify: true,
			},

			// Update fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_dedupe_processor" "my_processor" {
						title = "new title"
						description = "new desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						fields = [".thing3"]
						number_of_events = 4999
						comparison_type = "Ignore"
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_dedupe_processor.my_processor", map[string]any{
						"pipeline_id":      "#mezmo_pipeline.test_parent.id",
						"title":            "new title",
						"description":      "new desc",
						"generation_id":    "1",
						"inputs.#":         "1",
						"inputs.0":         "#mezmo_http_source.my_source.id",
						"number_of_events": "4999",
						"comparison_type":  "Ignore",
						"fields.#":         "1",
						"fields.0":         ".thing3",
					}),
				),
			},

			// Error: server-side validation
			{
				Config: GetCachedConfig(cacheKey) + `
				resource "mezmo_dedupe_processor" "my_processor" {
					pipeline_id = mezmo_pipeline.test_parent.id
					inputs = []
					fields = ["not-a-valid-field"]
				}`,
				ExpectError: regexp.MustCompile("be a valid data access syntax"),
			},
			// confirm manually deleted resources are recreated
			{
				Config: GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_dedupe_processor" "test_processor" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					inputs 			= []
					fields 			= [".thing1", ".thing2"]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_dedupe_processor.test_processor", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_dedupe_processor.test_processor", "title", "new title"),
					// delete the resource
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_dedupe_processor.test_processor",
					),
				),
				// verify resource will be re-created after refresh
				ExpectNonEmptyPlan: true,
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
