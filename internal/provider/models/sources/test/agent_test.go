package sources

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/providertest"
)

func TestAgentSourceResource(t *testing.T) {
        const cacheKey = "agent_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Error: pipeline_id is required
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}`) + `
					resource "mezmo_agent_source" "my_source" {
						itle = "my agent title"
						description = "my agent description"
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},
			// Create and Read testing
			{
			    Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_agent_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my agent title"
						description = "my agent description"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_agent_source.my_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestMatchResourceAttr(
						"mezmo_agent_source.my_source", "gateway_route_id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_agent_source.my_source", map[string]any{
						"description":      "my agent description",
						"generation_id":    "0",
						"title":            "my agent title",
						"capture_metadata": "false",
						"pipeline_id":      "#mezmo_pipeline.test_parent.id",
					}),
				),
			},
			// Update and Read testing
			{
			    Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_agent_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "new title"
						description = "new description"
						capture_metadata = "true"
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_agent_source.my_source", map[string]any{
						"description":      "new description",
						"generation_id":    "1",
						"title":            "new title",
						"capture_metadata": "true",
					}),
				),
			},
			// Supply gateway_route_id
			{
				Config: SetCachedConfig("Supply gateway_route_id", `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_agent_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my agent title"
						description = "my agent description"
					}`) + `
					resource "mezmo_agent_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "A shared agent source"
						description = "This source provides gateway_route_id"
						gateway_route_id = mezmo_agent_source.my_source.gateway_route_id
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_agent_source.my_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_agent_source.shared_source", map[string]any{
						"description":      "This source provides gateway_route_id",
						"generation_id":    "0",
						"title":            "A shared agent source",
						"capture_metadata": "false",
						"pipeline_id":      "#mezmo_pipeline.test_parent.id",
						"gateway_route_id": "#mezmo_agent_source.my_source.gateway_route_id",
					}),
				),
			},
			// Updating gateway_route_id is not allowed
			{
				Config: GetCachedConfig("Supply gateway_route_id") + `
					resource "mezmo_agent_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "A new title"
						gateway_route_id = mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile("This field is immutable after resource creation."),
			},
			// gateway_route_id can be specified if it's the same value
			{
				Config: GetCachedConfig("Supply gateway_route_id") + `
					resource "mezmo_agent_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "Updated title again"
						gateway_route_id = mezmo_agent_source.my_source.gateway_route_id
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_agent_source.my_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_agent_source.shared_source", map[string]any{
						"title":            "Updated title again",
						"gateway_route_id": "#mezmo_agent_source.my_source.gateway_route_id",
					}),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
