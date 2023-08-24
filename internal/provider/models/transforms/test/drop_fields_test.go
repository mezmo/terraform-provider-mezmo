package transforms

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/providertest"
)

func TestDropFieldsTransform(t *testing.T) {
	const cacheKey = "drop_fields_resources"
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
					resource "mezmo_drop_fields_transform" "my_transform" {
						fields = [".nope"]
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},

			// Error: `fields` is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_drop_fields_transform" "my_transform" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile("The argument \"fields\" is required"),
			},

			// Error: `fields` array length validation
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_drop_fields_transform" "my_transform" {
						pipeline_id = mezmo_pipeline.test_parent.id
						fields = []
					}`,
				ExpectError: regexp.MustCompile("Attribute fields list must contain at least 1 elements"),
			},

			// Error: `fields` values validates length
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_drop_fields_transform" "my_transform" {
						pipeline_id = mezmo_pipeline.test_parent.id
						fields = [""]
					}`,
				ExpectError: regexp.MustCompile("Attribute fields\\[0\\] string length must be at least 1"),
			},

			// Create
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_drop_fields_transform" "my_transform" {
						title = "transform title"
						description = "transform desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						fields = [".thing1", ".thing2"]
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_drop_fields_transform.my_transform", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_drop_fields_transform.my_transform", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "transform title",
						"description":   "transform desc",
						"generation_id": "0",
						"inputs.#":      "0",
						"fields.#":      "2",
						"fields.0":      ".thing1",
						"fields.1":      ".thing2",
					}),
				),
			},

			// Update fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_drop_fields_transform" "my_transform" {
						title = "new title"
						description = "new desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						fields = [".thing3"]
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_drop_fields_transform.my_transform", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "new title",
						"description":   "new desc",
						"generation_id": "1",
						"inputs.#":      "1",
						"inputs.0":      "#mezmo_http_source.my_source.id",
						"fields.#":      "1",
						"fields.0":      ".thing3",
					}),
				),
			},

			// Error: server-side validation
			{
				Config: GetCachedConfig(cacheKey) + `
				resource "mezmo_drop_fields_transform" "my_transform" {
					pipeline_id = mezmo_pipeline.test_parent.id
					inputs = []
					fields = ["not-a-valid-field"]
				}`,
				ExpectError: regexp.MustCompile("be a valid data access syntax"),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
