package processors

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/providertest"
)

func TestAccSampleProcessor(t *testing.T) {
	const cacheKey = "sample_resources"
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
					resource "mezmo_sample_processor" "my_processor" {
						field = ".nope"
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},

			// Error: `always_include` validation
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_sample_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						always_include = {

						}
					}`,
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"always_include\""),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_sample_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						always_include = {
							field = ".my_field"
							operator = "equal"
							value_number = "3587"
							case_sensitive = false
						}
					}`,
				ExpectError: regexp.MustCompile("(?s)Attribute \"always_include.value_string\" must be specified when.*\"always_include.case_sensitive\" is specified"),
			},

			// Create with negation
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_sample_processor" "my_new_processor" {
						title = "test title"
						description = "test desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						always_include = {
							field = ".my_field"
							operator = "exists"
							negate = true
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_sample_processor.my_new_processor", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_sample_processor.my_new_processor", map[string]any{
						"pipeline_id":             "#mezmo_pipeline.test_parent.id",
						"title":                   "test title",
						"description":             "test desc",
						"generation_id":           "0",
						"inputs.#":                "0",
						"rate":                    "10",
						"always_include.field":    ".my_field",
						"always_include.operator": "exists",
						"always_include.negate":   "true",
					}),
				),
			},

			// Create with defaults
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_sample_processor" "my_processor" {
						title = "test title"
						description = "test desc"
						pipeline_id = mezmo_pipeline.test_parent.id
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_sample_processor.my_processor", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_sample_processor.my_processor", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "test title",
						"description":   "test desc",
						"generation_id": "0",
						"inputs.#":      "0",
						"rate":          "10",
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_sample_processor" "import_target" {
						title = "test title"
						description = "test desc"
						pipeline_id = mezmo_pipeline.test_parent.id
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_sample_processor.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_sample_processor.my_processor"),
				ImportStateVerify: true,
			},

			// Update fields: always_include w/ no value
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_sample_processor" "my_processor" {
						title = "new title"
						description = "new desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						rate = 3444
						always_include = {
							field = ".my_field"
							operator = "exists"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_sample_processor.my_processor", map[string]any{
						"pipeline_id":             "#mezmo_pipeline.test_parent.id",
						"title":                   "new title",
						"description":             "new desc",
						"generation_id":           "1",
						"inputs.#":                "1",
						"inputs.0":                "#mezmo_http_source.my_source.id",
						"rate":                    "3444",
						"always_include.field":    ".my_field",
						"always_include.operator": "exists",
						"always_include.negate":   "false",
					}),
				),
			},

			// Update fields: add value_number
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_sample_processor" "my_processor" {
						title = "new title"
						description = "new desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						rate = 3444
						always_include = {
							field = ".my_field"
							operator = "greater"
							value_number = 122
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_sample_processor.my_processor", map[string]any{
						"pipeline_id":                 "#mezmo_pipeline.test_parent.id",
						"title":                       "new title",
						"description":                 "new desc",
						"generation_id":               "2",
						"inputs.#":                    "1",
						"inputs.0":                    "#mezmo_http_source.my_source.id",
						"rate":                        "3444",
						"always_include.field":        ".my_field",
						"always_include.operator":     "greater",
						"always_include.value_number": "122",
					}),
				),
			},

			// Update fields: add value_string
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_sample_processor" "my_processor" {
						title = "new title"
						description = "new desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						rate = 678
						always_include = {
							field = ".my_field"
							operator = "contains"
							value_string = "my text"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_sample_processor.my_processor", map[string]any{
						"pipeline_id":                   "#mezmo_pipeline.test_parent.id",
						"title":                         "new title",
						"description":                   "new desc",
						"generation_id":                 "3",
						"inputs.#":                      "1",
						"inputs.0":                      "#mezmo_http_source.my_source.id",
						"rate":                          "678",
						"always_include.field":          ".my_field",
						"always_include.operator":       "contains",
						"always_include.value_string":   "my text",
						"always_include.case_sensitive": "true",
					}),
				),
			},

			// Update fields: change case sensitivity
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_sample_processor" "my_processor" {
						title = "new title"
						description = "new desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						rate = 678
						always_include = {
							field = ".my_field"
							operator = "starts_with"
							value_string = "BEGIN"
							case_sensitive = false
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_sample_processor.my_processor", map[string]any{
						"pipeline_id":                   "#mezmo_pipeline.test_parent.id",
						"title":                         "new title",
						"description":                   "new desc",
						"generation_id":                 "4",
						"inputs.#":                      "1",
						"inputs.0":                      "#mezmo_http_source.my_source.id",
						"rate":                          "678",
						"always_include.field":          ".my_field",
						"always_include.operator":       "starts_with",
						"always_include.value_string":   "BEGIN",
						"always_include.case_sensitive": "false",
					}),
				),
			},

			// Error: server-side validation
			{
				Config: GetCachedConfig(cacheKey) + `
				resource "mezmo_sample_processor" "my_processor" {
					pipeline_id = mezmo_pipeline.test_parent.id
					inputs = []
					always_include = {
						field = ".my_field"
						operator = "greater"
						value_string = "my text"
					}
				}`,
				ExpectError: regexp.MustCompile("Value must be a"),
			},
			// confirm manually deleted resources are recreated
			{
				Config: GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_sample_processor" "test_processor" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					inputs 			= []
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_sample_processor.test_processor", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_sample_processor.test_processor", "title", "new title"),
					// delete the resource
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_sample_processor.test_processor",
					),
				),
				// verify resource will be re-created after refresh
				ExpectNonEmptyPlan: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
