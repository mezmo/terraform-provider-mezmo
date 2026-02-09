package processors

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/v5/internal/provider/providertest"
)

func TestAccMapFieldsProcessor(t *testing.T) {
	const cacheKey = "map_fields_reduce_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Error: pipeline_id is required
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_http_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`) + `
					resource "mezmo_map_fields_processor" "my_processor" {
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},

			// Error: Invalid field path
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_map_fields_processor" "single_mapping" {
						title = "some title"
						description = "some description"
						inputs = [mezmo_http_source.my_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id

						mappings = [
							{
								source_field = "invalid" // No preceding .
								target_field = ".field2"
							},
						]
					}`,
				ExpectError: regexp.MustCompile("(?s)source_field.*be a valid data access syntax"),
			},

			// Error: No mappings
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_map_fields_processor" "single_mapping" {
						title = "some title"
						description = "some description"
						inputs = [mezmo_http_source.my_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id

						mappings = []
					}`,
				ExpectError: regexp.MustCompile("(?s)Error: Invalid Attribute Value"),
			},

			// Error: Empty fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_map_fields_processor" "single_mapping" {
						title = "some title"
						description = "some description"
						inputs = [mezmo_http_source.my_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id

						mappings = [
							{
								source_field = "" // Empty
								target_field = "" // Empty
							}
						]
					}`,
				ExpectError: regexp.MustCompile("(?s)Error: Invalid Attribute Value"),
			},

			// Single mapping (default values)
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_map_fields_processor" "single_mapping" {
						title = "some title"
						description = "some description"
						inputs = [mezmo_http_source.my_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id

						mappings = [
							{
								source_field = ".field1"
								target_field = ".field2"
							},
						]
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_map_fields_processor.single_mapping", map[string]any{
						"pipeline_id":                 "#mezmo_pipeline.test_parent.id",
						"title":                       "some title",
						"description":                 "some description",
						"generation_id":               "0",
						"inputs.#":                    "1",
						"mappings.#":                  "1",
						"mappings.0.source_field":     ".field1",
						"mappings.0.target_field":     ".field2",
						"mappings.0.drop_source":      "false",
						"mappings.0.overwrite_target": "false",
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_map_fields_processor" "import_target" {
						title = "some title"
						description = "some description"
						inputs = [mezmo_http_source.my_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id

						mappings = [
							{
								source_field = ".field1"
								target_field = ".field2"
							},
						]
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_map_fields_processor.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_map_fields_processor.single_mapping"),
				ImportStateVerify: true,
			},

			// Multiple mappings
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_map_fields_processor" "multiple_mappings" {
						title = "some title"
						description = "some description"
						inputs = [mezmo_http_source.my_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id

						mappings = [
							{
								source_field = ".field1"
								target_field = ".field2"
							},
							{
								source_field = ".field3"
								target_field = ".field4"
								drop_source = true
							},
							{
								source_field = ".field5"
								target_field = ".field6"
								overwrite_target = true
							},
							{
								source_field = ".field7"
								target_field = ".field8"
								drop_source = false
								overwrite_target = true
							},
						]
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_map_fields_processor.multiple_mappings", map[string]any{
						"pipeline_id":                 "#mezmo_pipeline.test_parent.id",
						"title":                       "some title",
						"description":                 "some description",
						"generation_id":               "0",
						"inputs.#":                    "1",
						"mappings.#":                  "4",
						"mappings.0.source_field":     ".field1",
						"mappings.0.target_field":     ".field2",
						"mappings.0.drop_source":      "false",
						"mappings.0.overwrite_target": "false",
						"mappings.1.source_field":     ".field3",
						"mappings.1.target_field":     ".field4",
						"mappings.1.drop_source":      "true",
						"mappings.1.overwrite_target": "false",
						"mappings.2.source_field":     ".field5",
						"mappings.2.target_field":     ".field6",
						"mappings.2.drop_source":      "false",
						"mappings.2.overwrite_target": "true",
						"mappings.3.source_field":     ".field7",
						"mappings.3.target_field":     ".field8",
						"mappings.3.drop_source":      "false",
						"mappings.3.overwrite_target": "true",
					}),
				),
			},
			// confirm manually deleted resources are recreated
			{
				Config: GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_map_fields_processor" "test_processor" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					inputs 			= []
					mappings 		= [
						{
							source_field = ".field1"
							target_field = ".field2"
						},
					]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_map_fields_processor.test_processor", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_map_fields_processor.test_processor", "title", "new title"),
					// delete the resource
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_map_fields_processor.test_processor",
					),
				),
				// verify resource will be re-created after refresh
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
