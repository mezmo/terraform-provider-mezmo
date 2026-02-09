package sources

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/v5/internal/provider/providertest"
)

func TestAccLogAnalysisSource(t *testing.T) {
	cacheKey := "log_analysis_source_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Required fields parent pipeline id
			{
				Config: GetProviderConfig() + `
					resource "mezmo_log_analysis_source" "my_source" {
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},

			// Create and Read testing
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}`) + `
					resource "mezmo_log_analysis_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my source title"
						description = "my source description"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_log_analysis_source.my_source", "id", regexp.MustCompile(`[\w-]{36}`)),

					resource.TestCheckNoResourceAttr("mezmo_log_analysis_source.my_source", "shared_source_id"),

					StateHasExpectedValues("mezmo_log_analysis_source.my_source", map[string]any{
						"description":   "my source description",
						"generation_id": "0",
						"title":         "my source title",
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
					}),
					resource.TestCheckResourceAttrSet("mezmo_log_analysis_source.my_source", "generation_id"),
				),
			},
			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_log_analysis_source" "import_target" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my source title"
						description = "my source description"
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_log_analysis_source.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_log_analysis_source.my_source"),
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_log_analysis_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "new title"
						description = "new description"
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_log_analysis_source.my_source", map[string]any{
						"description":   "new description",
						"generation_id": "1",
						"title":         "new title",
					}),
				),
			},
			// confirm manually deleted resources are recreated
			{
				Config: GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_log_analysis_source" "test_source" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_log_analysis_source.test_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_log_analysis_source.test_source", "title", "new title"),
					// delete the resource
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_log_analysis_source.test_source",
					),
				),
				// verify resource will be re-created after refresh
				ExpectNonEmptyPlan: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
