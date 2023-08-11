package sinks

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/providertest"
)

func TestBlackholeSinkResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Required field pipeline id
			{
				Config: GetProviderConfig() + `
					resource "mezmo_blackhole_sink" "my_sink" {}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required, but no definition was found"),
			},
			// Create and Read testing
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_blackhole_sink" "my_sink" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my sink title"
						description = "my sink description"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_blackhole_sink.my_sink", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_blackhole_sink.my_sink", map[string]any{
						"description":   "my sink description",
						"generation_id": "0",
						"title":         "my sink title",
						"ack_enabled":   "true",
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
					resource "mezmo_blackhole_sink" "my_sink" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "new title"
						description = "new description"
						ack_enabled = false
						inputs = [mezmo_demo_source.my_source.id]
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_blackhole_sink.my_sink", map[string]any{
						"description":   "new description",
						"generation_id": "1",
						"title":         "new title",
						"ack_enabled":   "false",
						"inputs.#":      "1",
					}),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
