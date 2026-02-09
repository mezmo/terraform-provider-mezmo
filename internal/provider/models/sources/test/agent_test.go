package sources

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/v5/internal/provider/providertest"
)

func TestAccAgentSourceResource(t *testing.T) {
	const cacheKey = "agent_source_resources"
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
						"mezmo_agent_source.my_source", "shared_source_id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_agent_source.my_source", map[string]any{
						"description":   "my agent description",
						"generation_id": "0",
						"title":         "my agent title",
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
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
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_agent_source.my_source", map[string]any{
						"description":   "new description",
						"generation_id": "1",
						"title":         "new title",
					}),
				),
			},
			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_agent_source" "import_target" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "bad new title"
						description = "new description"
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_agent_source.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_agent_source.my_source"),
				ImportStateVerify: true,
			},
			// Supply shared_source_id
			{
				Config: SetCachedConfig(cacheKey, `
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
						description = "This source provides shared_source_id"
						shared_source_id = mezmo_agent_source.my_source.shared_source_id
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_agent_source.my_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_agent_source.shared_source", map[string]any{
						"description":      "This source provides shared_source_id",
						"generation_id":    "0",
						"title":            "A shared agent source",
						"pipeline_id":      "#mezmo_pipeline.test_parent.id",
						"shared_source_id": "#mezmo_agent_source.my_source.shared_source_id",
					}),
				),
			},
			// Updating shared_source_id is not allowed
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_agent_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "A new title"
						shared_source_id = mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile("This field is immutable after resource creation."),
			},
			// shared_source_id can be specified if it's the same value
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_agent_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "Updated title again"
						shared_source_id = mezmo_agent_source.my_source.shared_source_id
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_agent_source.my_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_agent_source.shared_source", map[string]any{
						"title":            "Updated title again",
						"shared_source_id": "#mezmo_agent_source.my_source.shared_source_id",
					}),
				),
			},
			// confirm manually deleted resources are recreated
			{
				Config: GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_agent_source" "test_source" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_agent_source.test_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_agent_source.test_source", "title", "new title"),
					// delete the resource
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_agent_source.test_source",
					),
				),
				// verify resource will be re-created after refresh
				ExpectNonEmptyPlan: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
