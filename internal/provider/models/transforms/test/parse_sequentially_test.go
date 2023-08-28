package transforms

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/providertest"
)

func TestParseSequentiallyTransform(t *testing.T) {
	const cacheKey = "dedupe_resources"
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
					resource "mezmo_parse_sequentially_transform" "my_transform" {
						fields = [".nope"]
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},

			// Error: `field` is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_transform" "my_transform" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile("The argument \"field\" is required"),
			},

			// Error: `parsers` is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_transform" "my_transform" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
					}`,
				ExpectError: regexp.MustCompile("The argument \"parsers\" is required"),
			},

			// Error: `parsers` array length validation
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_transform" "my_transform" {
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parsers = []
					}`,
				ExpectError: regexp.MustCompile("Attribute parsers list must contain at least 1 elements"),
			},

			// Create with defaults
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_parse_sequentially_transform" "my_transform" {
						title = "title"
						description = "desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						field = ".something"
						parsers = [
							{
								parser = "parse_csv"
								label = "my label"
								options = {
									first = "yass"
									second = "nooo"
								}
							}
						]
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_parse_sequentially_transform.my_transform", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_parse_sequentially_transform.my_transform", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "title",
						"description":   "desc",
						"generation_id": "0",
						"inputs.#":      "0",
					}),
				),
			},

			// // Update fields
			// {
			// 	Config: GetCachedConfig(cacheKey) + `
			// 		resource "mezmo_parse_sequentially_transform" "my_transform" {
			// 			title = "new title"
			// 			description = "new desc"
			// 			pipeline_id = mezmo_pipeline.test_parent.id
			// 			inputs = [mezmo_http_source.my_source.id]
			// 			fields = [".thing3"]
			// 			number_of_events = 4999
			// 			comparison_type = "Ignore"
			// 		}`,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		StateHasExpectedValues("mezmo_parse_sequentially_transform.my_transform", map[string]any{
			// 			"pipeline_id":      "#mezmo_pipeline.test_parent.id",
			// 			"title":            "new title",
			// 			"description":      "new desc",
			// 			"generation_id":    "1",
			// 			"inputs.#":         "1",
			// 			"inputs.0":         "#mezmo_http_source.my_source.id",
			// 			"number_of_events": "4999",
			// 			"comparison_type":  "Ignore",
			// 			"fields.#":         "1",
			// 			"fields.0":         ".thing3",
			// 		}),
			// 	),
			// },

			// // Error: server-side validation
			// {
			// 	Config: GetCachedConfig(cacheKey) + `
			// 	resource "mezmo_parse_sequentially_transform" "my_transform" {
			// 		pipeline_id = mezmo_pipeline.test_parent.id
			// 		inputs = []
			// 		fields = ["not-a-valid-field"]
			// 	}`,
			// 	ExpectError: regexp.MustCompile("be a valid data access syntax"),
			// },

			// Delete testing automatically occurs in TestCase
		},
	})
}
