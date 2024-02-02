package sources

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestDemoSourceResource(t *testing.T) {
	const cacheKey = "demo_source_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Required fields json
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_demo_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			// Required fields parent pipeline id
			{
				Config: GetProviderConfig() + `
					resource "mezmo_demo_source" "my_source" {
						format = "json"
					}`,
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			// Required fields parent pipeline id
			{
				Config: GetProviderConfig() + `
					resource "mezmo_demo_source" "my_source" {
						pipeline_id = "798e1028-0b60-11ee-be56-0242ac120002"
						format = "NOT_VALID"
					}`,
				ExpectError: regexp.MustCompile("Attribute format value must be one of"),
			},
			// Create and Read testing
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}`) + `
					resource "mezmo_demo_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my source title"
						description = "my source description"
						format = "json"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_demo_source.my_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_demo_source.my_source", map[string]any{
						"description":      "my source description",
						"generation_id":    "0",
						"title":            "my source title",
						"format":           "json",
						"pipeline_id":      "#mezmo_pipeline.test_parent.id",
						"gateway_route_id": nil,
					}),
					resource.TestCheckResourceAttrSet("mezmo_demo_source.my_source", "generation_id"),
				),
			},
			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_demo_source" "import_target" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my source title"
						description = "my source description"
						format = "json"
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_demo_source.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_demo_source.my_source"),
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_demo_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "new title"
						description = "new description"
						format = "apache_common"
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_demo_source.my_source", map[string]any{
						"description":   "new description",
						"generation_id": "1",
						"title":         "new title",
						"format":        "apache_common",
					}),
				),
			},
			// confirm manually deleted resources are recreated
			{
				Config: GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_demo_source" "test_source" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					format 			= "json"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_demo_source.test_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_demo_source.test_source", "title", "new title"),
					// delete the resource
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_demo_source.test_source",
					),
				),
				// verify resource will be re-created after refresh
				ExpectNonEmptyPlan: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
