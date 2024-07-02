package processors

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestAccFlattenFieldsProcessor(t *testing.T) {
	const cacheKey = "flatten_fields_resources"
	SetCachedConfig(cacheKey, `
		resource "mezmo_pipeline" "test_parent" {
			title = "pipeline"
		}
		resource "mezmo_http_source" "my_source" {
			pipeline_id = mezmo_pipeline.test_parent.id
		}`,
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Error: pipeline_id is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_flatten_fields_processor" "my_processor" {
						fields = [".nope"]
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},

			// Error: `fields` values validates length
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_flatten_fields_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						fields = [""]
					}`,
				ExpectError: regexp.MustCompile(`Attribute fields\[0\] string length must be at least 1`),
			},

			// Create with defaults
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_flatten_fields_processor" "my_processor_defaults" {
						title = "flatten fields title"
						description = "flatten fields desc"
						pipeline_id = mezmo_pipeline.test_parent.id
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_flatten_fields_processor.my_processor_defaults", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_flatten_fields_processor.my_processor_defaults", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "flatten fields title",
						"description":   "flatten fields desc",
						"generation_id": "0",
						"inputs.#":      "0",
						"delimiter":     "_",
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_flatten_fields_processor" "import_target" {
						title = "flatten fields title"
						description = "flatten fields desc"
						pipeline_id = mezmo_pipeline.test_parent.id
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_flatten_fields_processor.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_flatten_fields_processor.my_processor_defaults"),
				ImportStateVerify: true,
			},

			// Create with fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_flatten_fields_processor" "my_processor" {
						title = "flatten fields title"
						description = "flatten fields desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						fields = [".thing1", ".thing2"]
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_flatten_fields_processor.my_processor", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_flatten_fields_processor.my_processor", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "flatten fields title",
						"description":   "flatten fields desc",
						"generation_id": "0",
						"inputs.#":      "0",
						"delimiter":     "_",
						"fields.#":      "2",
						"fields.0":      ".thing1",
						"fields.1":      ".thing2",
					}),
				),
			},

			// Update fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_flatten_fields_processor" "my_processor" {
						title = "new title"
						description = "new desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						fields = [".thing3"]
						delimiter = "~"
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_flatten_fields_processor.my_processor", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "new title",
						"description":   "new desc",
						"generation_id": "1",
						"inputs.#":      "1",
						"inputs.0":      "#mezmo_http_source.my_source.id",
						"delimiter":     "~",
						"fields.#":      "1",
						"fields.0":      ".thing3",
					}),
				),
			},

			// Update to unset fields
			{
				Config: GetCachedConfig(cacheKey) + `
				resource "mezmo_flatten_fields_processor" "my_processor" {
					title = "new title"
					description = "new desc"
					pipeline_id = mezmo_pipeline.test_parent.id
					inputs = [mezmo_http_source.my_source.id]
					fields = null
					delimiter = "~"
				}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_flatten_fields_processor.my_processor", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "new title",
						"description":   "new desc",
						"generation_id": "2",
						"inputs.#":      "1",
						"inputs.0":      "#mezmo_http_source.my_source.id",
						"delimiter":     "~",
					}),
				),
			},

			// Error: server-side validation
			{
				Config: GetCachedConfig(cacheKey) + `
				resource "mezmo_flatten_fields_processor" "my_processor" {
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
				resource "mezmo_flatten_fields_processor" "test_processor" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					inputs 			= []
					fields 			= [".thing1", ".thing2"]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_flatten_fields_processor.test_processor", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_flatten_fields_processor.test_processor", "title", "new title"),
					// delete the resource
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_flatten_fields_processor.test_processor",
					),
				),
				// verify resource will be re-created after refresh
				ExpectNonEmptyPlan: true,
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
