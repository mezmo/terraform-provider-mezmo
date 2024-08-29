package sources

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/providertest"
)

func TestAccKinesisFirehoseSourceResource(t *testing.T) {
	cacheKey := "kinesis_firehose_source_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { providertest.TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Required field test cases
			{
				Config: providertest.GetProviderConfig() + `
					resource "mezmo_kinesis_firehose_source" "req_field_source" {}`,
				ExpectError: regexp.MustCompile("argument \"pipeline_id\" is required"),
			},
			// Validator test cases
			{
				Config: providertest.GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_kinesis_firehose_source" "val_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						decoding = "invalid"
					}`,
				ExpectError: regexp.MustCompile("Attribute decoding value must be one of:"),
			},
			// Create and Update State
			{
				Config: providertest.SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}`) + `
					resource "mezmo_kinesis_firehose_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "test title"
						description = "test description"
						decoding = "text"
						capture_metadata = true
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_kinesis_firehose_source.my_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					providertest.StateHasExpectedValues("mezmo_kinesis_firehose_source.my_source", map[string]any{
						"title":            "test title",
						"description":      "test description",
						"generation_id":    "0",
						"decoding":         "text",
						"capture_metadata": "true",
					}),
				),
			},
			// Import
			{
				Config: providertest.GetCachedConfig(cacheKey) + `
					resource "mezmo_kinesis_firehose_source" "import_target" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "test title"
						description = "test description"
						decoding = "text"
						capture_metadata = true
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_kinesis_firehose_source.import_target",
				ImportStateIdFunc: providertest.ComputeImportId("mezmo_kinesis_firehose_source.my_source"),
				ImportStateVerify: true,
			},
			{
				Config: providertest.GetCachedConfig(cacheKey) + `
					resource "mezmo_kinesis_firehose_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "new title"
						description = "new description"
						decoding = "json"
						capture_metadata = false
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					providertest.StateHasExpectedValues("mezmo_kinesis_firehose_source.my_source", map[string]any{
						"title":            "new title",
						"description":      "new description",
						"generation_id":    "1",
						"decoding":         "json",
						"capture_metadata": "false",
					}),
				),
			},
			// Supply shared_source_id
			{
				Config: providertest.SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_kinesis_firehose_source" "parent_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "parent"
						description = "parent kinesis source"
					}`) + `
					resource "mezmo_kinesis_firehose_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "shared"
						description = "shared kinesis source"
						shared_source_id = mezmo_kinesis_firehose_source.parent_source.id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_kinesis_firehose_source.shared_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					providertest.StateHasExpectedValues("mezmo_kinesis_firehose_source.shared_source", map[string]any{
						"title":            "shared",
						"description":      "shared kinesis source",
						"generation_id":    "0",
						"decoding":         "json",
						"capture_metadata": "false",
						"pipeline_id":      "#mezmo_pipeline.test_parent.id",
						"shared_source_id": "#mezmo_kinesis_firehose_source.parent_source.shared_source_id",
					}),
				),
			},
			// Updating shared_source_id is not allowed
			{
				Config: providertest.GetCachedConfig(cacheKey) + `
					resource "mezmo_kinesis_firehose_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						shared_source_id =  mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile("This field is immutable after resource creation."),
			},
			// shared_source_id can be specified if it's the same value
			{
				Config: providertest.GetCachedConfig(cacheKey) + `
					resource "mezmo_kinesis_firehose_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "another title update"
						shared_source_id = mezmo_kinesis_firehose_source.parent_source.shared_source_id
					}`,
				Check: resource.ComposeTestCheckFunc(
					providertest.StateHasExpectedValues("mezmo_kinesis_firehose_source.shared_source", map[string]any{
						"title":            "another title update",
						"generation_id":    "1",
						"shared_source_id": "#mezmo_kinesis_firehose_source.parent_source.shared_source_id",
					}),
				),
			},
			// confirm manually deleted resources are recreated
			{
				Config: providertest.GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_kinesis_firehose_source" "test_source" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_kinesis_firehose_source.test_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_kinesis_firehose_source.test_source", "title", "new title"),
					// delete the resource
					providertest.TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_kinesis_firehose_source.test_source",
					),
				),
				// verify resource will be re-created after refresh
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
