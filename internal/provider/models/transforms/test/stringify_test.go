package transforms

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/providertest"
)

func TestStringifyTransformResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Required field pipeline id
			{
				Config: GetProviderConfig() + `
					resource "mezmo_stringify_transform" "my_transform" {}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required, but no definition was found"),
			},
			// Create and Read testing
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_stringify_transform" "my_transform" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my stringify title"
						description = "my stringify description"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_stringify_transform.my_transform", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_stringify_transform.my_transform", map[string]any{
						"description":   "my stringify description",
						"generation_id": "0",
						"title":         "my stringify title",
						"inputs.#":      "0",
					}),
				),
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
					resource "mezmo_stringify_transform" "my_transform" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "new title"
						description = "new description"
						inputs = [mezmo_demo_source.my_source.id]
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_stringify_transform.my_transform", map[string]any{
						"description":   "new description",
						"generation_id": "1",
						"title":         "new title",
						"inputs.#":      "1",
						"inputs.0":      "#mezmo_demo_source.my_source.id",
					}),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
