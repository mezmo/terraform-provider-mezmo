package sources

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestSplunkHecSource(t *testing.T) {
	cacheKey := "splunk_hec_source_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Required fields parent pipeline id
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}`) + `
					resource "mezmo_splunk_hec_source" "my_source" {}
				`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},
			// Create and Read testing
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_splunk_hec_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my title"
						description = "my description"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_splunk_hec_source.my_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestMatchResourceAttr(
						"mezmo_splunk_hec_source.my_source", "gateway_route_id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_splunk_hec_source.my_source", map[string]any{
						"description":      "my description",
						"title":            "my title",
						"generation_id":    "0",
						"capture_metadata": "false",
						"pipeline_id":      "#mezmo_pipeline.test_parent.id",
					}),
				),
			},
			// Update and Read testing
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_splunk_hec_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "new title"
						description = "new description"
						capture_metadata = "true"
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_splunk_hec_source.my_source", map[string]any{
						"description":      "new description",
						"generation_id":    "1",
						"title":            "new title",
						"capture_metadata": "true",
					}),
				),
			},
			// Supply gateway_route_id
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_splunk_hec_source" "parent_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my parent title"
						description = "my parent description"
					}`) + `
					resource "mezmo_splunk_hec_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "A shared splunk HEC source"
						description = "This source provides gateway_route_id"
						gateway_route_id = mezmo_splunk_hec_source.parent_source.gateway_route_id
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_splunk_hec_source.shared_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_splunk_hec_source.shared_source", map[string]any{
						"description":      "This source provides gateway_route_id",
						"title":            "A shared splunk HEC source",
						"generation_id":    "0",
						"capture_metadata": "false",
						"pipeline_id":      "#mezmo_pipeline.test_parent.id",
						"gateway_route_id": "#mezmo_splunk_hec_source.parent_source.gateway_route_id",
					}),
				),
			},
			// Updating gateway_route_id is not allowed
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_splunk_hec_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						gateway_route_id = mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile("This field is immutable after resource creation."),
			},
			// gateway_route_id can be specified if it's the same value
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_splunk_hec_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "Updated title"
						gateway_route_id = mezmo_splunk_hec_source.parent_source.gateway_route_id
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_splunk_hec_source.shared_source", map[string]any{
						"title":            "Updated title",
						"gateway_route_id": "#mezmo_splunk_hec_source.parent_source.gateway_route_id",
					}),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
