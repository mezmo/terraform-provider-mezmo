package sources

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/providertest"
)

func TestFluentSource(t *testing.T) {
	cacheKey := "fluent_source_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Error: requires "pipeline_id"
			{
				Config: GetProviderConfig() + `
					resource "mezmo_fluent_source" "my_source" {}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},

			// Error: "decoding" is an enum
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}`) + `
					resource "mezmo_fluent_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						decoding = "nope"
					}`,
				ExpectError: regexp.MustCompile("Attribute decoding value must be one of:"),
			},

			// Create with defaults
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_fluent_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my title"
						description = "my description"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_fluent_source.my_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestMatchResourceAttr(
						"mezmo_fluent_source.my_source", "gateway_route_id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_fluent_source.my_source", map[string]any{
						"description":      "my description",
						"title":            "my title",
						"pipeline_id":      "#mezmo_pipeline.test_parent.id",
						"generation_id":    "0",
						"decoding":         "json",
						"capture_metadata": "false",
					}),
				),
			},

			// Update and Read testing
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_fluent_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "new title"
						description = "new description"
						decoding = "ndjson"
						capture_metadata = "true"
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_fluent_source.my_source", map[string]any{
						"description":      "new description",
						"title":            "new title",
						"generation_id":    "1",
						"decoding":         "ndjson",
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
					resource "mezmo_http_source" "parent_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my http title"
						description = "my http description"
					}`) + `
					resource "mezmo_fluent_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "A shared source"
						description = "This source provides gateway_route_id"
						gateway_route_id = mezmo_http_source.parent_source.gateway_route_id
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_fluent_source.shared_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_fluent_source.shared_source", map[string]any{
						"description":      "This source provides gateway_route_id",
						"title":            "A shared source",
						"pipeline_id":      "#mezmo_pipeline.test_parent.id",
						"generation_id":    "0",
						"decoding":         "json",
						"capture_metadata": "false",
						"gateway_route_id": "#mezmo_http_source.parent_source.gateway_route_id",
					}),
				),
			},

			// Updating gateway_route_id is not allowed
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_fluent_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						gateway_route_id = mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile("This field is immutable after resource creation."),
			},

			// gateway_route_id can be specified if it's the same value
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_fluent_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "Updated title"
						gateway_route_id = mezmo_http_source.parent_source.gateway_route_id
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_fluent_source.shared_source", map[string]any{
						"title":            "Updated title",
						"generation_id":    "1",
						"gateway_route_id": "#mezmo_http_source.parent_source.gateway_route_id",
					}),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
