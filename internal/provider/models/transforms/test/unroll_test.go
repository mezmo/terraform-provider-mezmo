package transforms

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/providertest"
)

func TestUnrollTransform(t *testing.T) {
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
					resource "mezmo_unroll_transform" "my_transform" {
						field = ".nope"
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},

			// Error: `field` is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_unroll_transform" "my_transform" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile("The argument \"field\" is required"),
			},

			// Error: `field` values validates length
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_unroll_transform" "my_transform" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ""
					}`,
				ExpectError: regexp.MustCompile("Attribute field string length must be at least 1"),
			},

			// Create
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_unroll_transform" "my_transform" {
						title = "transform title"
						description = "transform desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".thing1"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_unroll_transform.my_transform", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_unroll_transform.my_transform", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "transform title",
						"description":   "transform desc",
						"generation_id": "0",
						"inputs.#":      "0",
						"field":         ".thing1",
						"values_only":   "true",
					}),
				),
			},

			// Update
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_unroll_transform" "my_transform" {
						title = "new title"
						description = "new desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						field = ".thing2"
						values_only = false
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_unroll_transform.my_transform", map[string]any{
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
				resource "mezmo_unroll_transform" "my_transform" {
					pipeline_id = mezmo_pipeline.test_parent.id
					inputs = []
					field = "not-a-valid-field"
				}`,
				ExpectError: regexp.MustCompile("match pattern"),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
