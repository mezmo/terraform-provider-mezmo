package sources

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/providertest"
)

func TestAccFluentSource(t *testing.T) {
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
						"mezmo_fluent_source.my_source", "shared_source_id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_fluent_source.my_source", map[string]any{
						"description":   "my description",
						"title":         "my title",
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"generation_id": "0",
						"decoding":      "json",
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_fluent_source" "import_target" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my title"
						description = "my description"
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_fluent_source.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_fluent_source.my_source"),
				ImportStateVerify: true,
			},

			// Update and Read testing
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_fluent_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "new title"
						description = "new description"
						decoding = "ndjson"
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_fluent_source.my_source", map[string]any{
						"description":   "new description",
						"title":         "new title",
						"generation_id": "1",
						"decoding":      "ndjson",
					}),
				),
			},

			// Supply shared_source_id
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_fluent_source" "parent_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my http title"
						description = "my http description"
					}`) + `
					resource "mezmo_fluent_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "A shared source"
						description = "This source provides shared_source_id"
						shared_source_id = mezmo_fluent_source.parent_source.shared_source_id
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_fluent_source.shared_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_fluent_source.shared_source", map[string]any{
						"description":      "This source provides shared_source_id",
						"title":            "A shared source",
						"pipeline_id":      "#mezmo_pipeline.test_parent.id",
						"generation_id":    "0",
						"decoding":         "json",
						"shared_source_id": "#mezmo_fluent_source.parent_source.shared_source_id",
					}),
				),
			},

			// Updating shared_source_id is not allowed
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_fluent_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						shared_source_id = mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile("This field is immutable after resource creation."),
			},

			// shared_source_id can be specified if it's the same value
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_fluent_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "Updated title"
						shared_source_id = mezmo_fluent_source.parent_source.shared_source_id
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_fluent_source.shared_source", map[string]any{
						"title":            "Updated title",
						"generation_id":    "1",
						"shared_source_id": "#mezmo_fluent_source.parent_source.shared_source_id",
					}),
				),
			},
			// confirm manually deleted resources are recreated
			{
				Config: GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_fluent_source" "test_source" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_fluent_source.test_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_fluent_source.test_source", "title", "new title"),
					// delete the resource
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_fluent_source.test_source",
					),
				),
				// verify resource will be re-created after refresh
				ExpectNonEmptyPlan: true,
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
