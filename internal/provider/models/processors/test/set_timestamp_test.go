package processors

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/providertest"
)

func TestAccSetTimestampProcessor_error(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Errors
			{
				Config: GetProviderConfig() + `
					resource "mezmo_set_timestamp_processor" "my_processor" {
						title = "some title"
							description = "some description"
							inputs = ["abc"]
							pipeline_id = "pipeline-id"
					}`,
				ExpectError: regexp.MustCompile("The argument \"parsers\" is required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_set_timestamp_processor" "my_processor" {
						title = "some title"
							description = "some description"
							inputs = ["abc"]
							pipeline_id = "pipeline-id"
							parsers = []
					}`,
				ExpectError: regexp.MustCompile("Attribute parsers list must contain at least 1 elements"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_set_timestamp_processor" "my_processor" {
						title = "some title"
							description = "some description"
							inputs = ["abc"]
							pipeline_id = "pipeline-id"
							parsers = [
								{
									timestamp_format = "%Y-%m-%dT%H:%M:%S"
								}
							]
					}`,
				ExpectError: regexp.MustCompile("\"field\" is\\s*required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_set_timestamp_processor" "my_processor" {
						title = "some title"
							description = "some description"
							inputs = ["abc"]
							pipeline_id = "pipeline-id"
							parsers = [
								{
									field = ".field"
								}
							]
					}`,
				ExpectError: regexp.MustCompile("\"timestamp_format\" is required"),
			},
		},
	})
}

func TestAccSetTimestampProcessor_validators(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		ErrorCheck: CheckMultipleErrors([]string{
			"field string length must be at least 1",
			"timestamp_format string length must be at least 1",
		}),
		Steps: []resource.TestStep{
			{
				Config: GetProviderConfig() + `
					resource "mezmo_set_timestamp_processor" "my_processor" {
						title = "some title"
							description = "some description"
							inputs = ["abc"]
							pipeline_id = "pipeline-id"
							parsers = [
								{
									field = ""
									timestamp_format = ""
								}
							]
					}`,
			},
		},
	})
}

func TestAccSetTimestampProcessor_crud(t *testing.T) {
	const cacheKey = "set_timestamp_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			{
				// create
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_http_source" "test_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`) + `
					resource "mezmo_set_timestamp_processor" "crud" {
						title = "some title"
						description = "some description"
						inputs = [mezmo_http_source.test_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id

						parsers = [
							{
								field = ".field1"
								timestamp_format = "%Y-%m-%dT%H:%M:%S"
							},
						]
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_set_timestamp_processor.crud", map[string]any{
						"pipeline_id":                "#mezmo_pipeline.test_parent.id",
						"title":                      "some title",
						"description":                "some description",
						"generation_id":              "0",
						"inputs.#":                   "1",
						"parsers.#":                  "1",
						"parsers.0.field":            ".field1",
						"parsers.0.timestamp_format": "%Y-%m-%dT%H:%M:%S",
					}),
				),
			},
			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_set_timestamp_processor" "import_target" {
						title = "some title"
						description = "some description"
						inputs = [mezmo_http_source.test_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id

						parsers = [
							{
								field = ".field1"
								timestamp_format = "%Y-%m-%dT%H:%M:%S"
							},
						]
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_set_timestamp_processor.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_set_timestamp_processor.crud"),
				ImportStateVerify: true,
			},
			// Update
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_http_source" "new_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}
					resource "mezmo_set_timestamp_processor" "crud" {
						title = "some title 2"
						description = "some description 2"
						inputs = [mezmo_http_source.new_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id

						parsers = [
							{
								field = ".field2"
								timestamp_format = "%Y-%m-%dT%H:%M:%S"
							},
							{
								field = ".field3"
								timestamp_format = "%m/%d/%Y::%H:%M"
							},
						]
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_set_timestamp_processor.crud", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_set_timestamp_processor.crud", map[string]any{
						"pipeline_id":                "#mezmo_pipeline.test_parent.id",
						"title":                      "some title 2",
						"description":                "some description 2",
						"generation_id":              "1",
						"inputs.#":                   "1",
						"inputs.0":                   "#mezmo_http_source.new_source.id",
						"parsers.#":                  "2",
						"parsers.0.field":            ".field2",
						"parsers.0.timestamp_format": "%Y-%m-%dT%H:%M:%S",
						"parsers.1.field":            ".field3",
						"parsers.1.timestamp_format": "%m/%d/%Y::%H:%M",
					}),
				),
			},
			// confirm manually deleted resources are recreated
			{
				Config: GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_http_source" "my_source2" {
					pipeline_id = mezmo_pipeline.test_parent2.id
				}
				resource "mezmo_set_timestamp_processor" "test_processor" {
						title = "new title"
						description = "new description"
						inputs = [mezmo_http_source.my_source2.id]
						pipeline_id = mezmo_pipeline.test_parent2.id

						parsers = [
							{
								field = ".field1"
								timestamp_format = "%Y-%m-%dT%H:%M:%S"
							},
						]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_set_timestamp_processor.test_processor", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_set_timestamp_processor.test_processor", "title", "new title"),
					// verify resource will be re-created after refresh
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_set_timestamp_processor.test_processor",
					),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
