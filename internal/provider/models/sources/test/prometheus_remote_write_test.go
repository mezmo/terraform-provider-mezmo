package sources

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/providertest"
)

func TestAccPrometheusRemoteWriteSource(t *testing.T) {
	const cacheKey = "prometheus_remote_write_source_resources"
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
					resource "mezmo_prometheus_remote_write_source" "my_source" {
						title = "my kafka title"
						description = "my kafka description"
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},
			// Create and Read testing
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}`) + `
					resource "mezmo_prometheus_remote_write_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my prometheus remote write title"
						description = "my prometheus remote write description"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_prometheus_remote_write_source.my_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestMatchResourceAttr(
						"mezmo_prometheus_remote_write_source.my_source", "shared_source_id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_prometheus_remote_write_source.my_source", map[string]any{
						"description":      "my prometheus remote write description",
						"generation_id":    "0",
						"title":            "my prometheus remote write title",
						"capture_metadata": "false",
						"pipeline_id":      "#mezmo_pipeline.test_parent.id",
					}),
				),
			},
			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_prometheus_remote_write_source" "import_target" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my prometheus remote write title"
						description = "my prometheus remote write description"
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_prometheus_remote_write_source.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_prometheus_remote_write_source.my_source"),
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_prometheus_remote_write_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "new title"
						description = "new description"
						capture_metadata = "true"
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_prometheus_remote_write_source.my_source", map[string]any{
						"description":      "new description",
						"generation_id":    "1",
						"title":            "new title",
						"capture_metadata": "true",
					}),
				),
			},
			// Supply shared_source_id
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_prometheus_remote_write_source" "parent_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my prometheus remote write title"
						description = "my prometheus remote write description"
					}`) + `
					resource "mezmo_prometheus_remote_write_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "A shared prometheus remote write source"
						description = "This source provides shared_source_id"
						shared_source_id = mezmo_prometheus_remote_write_source.parent_source.shared_source_id
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_prometheus_remote_write_source.shared_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_prometheus_remote_write_source.shared_source", map[string]any{
						"description":      "This source provides shared_source_id",
						"generation_id":    "0",
						"title":            "A shared prometheus remote write source",
						"capture_metadata": "false",
						"pipeline_id":      "#mezmo_pipeline.test_parent.id",
						"shared_source_id": "#mezmo_prometheus_remote_write_source.parent_source.shared_source_id",
					}),
				),
			},
			// Updating shared_source_id is not allowed
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_prometheus_remote_write_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						shared_source_id = mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile("This field is immutable after resource creation."),
			},
			// shared_source_id can be specified if it's the same value
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_prometheus_remote_write_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "Updated title"
						shared_source_id = mezmo_prometheus_remote_write_source.parent_source.shared_source_id
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_prometheus_remote_write_source.shared_source", map[string]any{
						"title":            "Updated title",
						"generation_id":    "1",
						"shared_source_id": "#mezmo_prometheus_remote_write_source.parent_source.shared_source_id",
					}),
				),
			},
			// confirm manually deleted resources are recreated
			{
				Config: GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_prometheus_remote_write_source" "test_source" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_prometheus_remote_write_source.test_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_prometheus_remote_write_source.test_source", "title", "new title"),
					// delete the resource
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_prometheus_remote_write_source.test_source",
					),
				),
				// verify resource will be re-created after refresh
				ExpectNonEmptyPlan: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
