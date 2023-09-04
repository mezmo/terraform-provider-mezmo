package processors

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/providertest"
)

func TestSampleProcessor(t *testing.T) {
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
						"pipeline_id":                 "#mezmo_pipeline.test_parent.id",
						"title":                       "new title",
						"description":                 "new desc",
						"generation_id":               "3",
						"inputs.#":                    "1",
						"inputs.0":                    "#mezmo_http_source.my_source.id",
						"rate":                        "678",
						"always_include.field":        ".my_field",
						"always_include.operator":     "contains",
						"always_include.value_string": "my text",
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

			// Delete testing automatically occurs in TestCase
		},
	})
}
