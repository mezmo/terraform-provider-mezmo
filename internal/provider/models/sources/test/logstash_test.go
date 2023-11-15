package sources

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestLogStashSource(t *testing.T) {
	cacheKey := "logstash_source_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Error: Required field "pipeline_id"
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
							title = "parent pipeline"
						}`) + `
					resource "mezmo_logstash_source" "my_source" {
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},
			// Error: "format" is an enum
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_logstash_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						format = "nope"
					}`,
				ExpectError: regexp.MustCompile("Attribute format value must be one of:"),
			},

			// Create and Read testing with defaults
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_logstash_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						description = "my description"
						title = "my title"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("mezmo_logstash_source.my_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestMatchResourceAttr("mezmo_logstash_source.my_source", "gateway_route_id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_logstash_source.my_source", map[string]any{
						"description":      "my description",
						"title":            "my title",
						"generation_id":    "0",
						"format":           "json",
						"capture_metadata": "false",
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_logstash_source" "import_target" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "bad new title"
						description = "new description"
						capture_metadata = "true"
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_logstash_source.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_logstash_source.my_source"),
				ImportStateVerify: true,
			},

			// Updates
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_logstash_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						description = "new description"
						title = "new title"
						format = "text"
						capture_metadata = "true"
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_logstash_source.my_source", map[string]any{
						"description":      "new description",
						"title":            "new title",
						"generation_id":    "1",
						"format":           "text",
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
					resource "mezmo_logstash_source" "parent_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my title"
						description = "my description"
					}`) + `
					resource "mezmo_logstash_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "A shared source"
						description = "This source provides gateway_route_id"
						gateway_route_id = mezmo_logstash_source.parent_source.gateway_route_id
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_logstash_source.shared_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_logstash_source.shared_source", map[string]any{
						"description":      "This source provides gateway_route_id",
						"generation_id":    "0",
						"title":            "A shared source",
						"format":           "json",
						"capture_metadata": "false",
						"pipeline_id":      "#mezmo_pipeline.test_parent.id",
						"gateway_route_id": "#mezmo_logstash_source.parent_source.gateway_route_id",
					}),
				),
			},

			// Updating gateway_route_id is not allowed
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_logstash_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						gateway_route_id = mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile("This field is immutable after resource creation."),
			},

			// gateway_route_id can be specified if it's the same value
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_logstash_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "Updated title"
						gateway_route_id = mezmo_logstash_source.parent_source.gateway_route_id
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_logstash_source.shared_source", map[string]any{
						"title":            "Updated title",
						"generation_id":    "1",
						"gateway_route_id": "#mezmo_logstash_source.parent_source.gateway_route_id",
					}),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
