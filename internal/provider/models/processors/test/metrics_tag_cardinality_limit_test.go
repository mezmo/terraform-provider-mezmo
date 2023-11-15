package processors

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestMetricsTagCardinalityLimitProcessor(t *testing.T) {
	const cacheKey = "metrics_tag_cardinality_limit_resources"
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
					resource "mezmo_metrics_tag_cardinality_limit_processor" "my_processor" {
						value_limit = 50
					}`,
				ExpectError: regexp.MustCompile(`The argument "pipeline_id" is required`),
			},

			// Error: value_limit is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_metrics_tag_cardinality_limit_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile(`The argument "value_limit" is required`),
			},

			// Error: length validation for tags
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_metrics_tag_cardinality_limit_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						value_limit = 50
						tags = [""]
					}`,
				ExpectError: regexp.MustCompile("Attribute tags.* string length must be at least 1"),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_metrics_tag_cardinality_limit_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						value_limit = 50
						tags = ["` + strings.Repeat("x", 101) + `"]
					}`,
				ExpectError: regexp.MustCompile("Attribute tags.* string length must be at most 100"),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_metrics_tag_cardinality_limit_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						value_limit = 50
						tags = ["x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x"]
					}`,
				ExpectError: regexp.MustCompile("Attribute tags list must contain at most 10 elements"),
			},

			// Error: length validation for exclude_tags
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_metrics_tag_cardinality_limit_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						value_limit = 50
						exclude_tags = [""]
					}`,
				ExpectError: regexp.MustCompile("Attribute exclude_tags.* string length must be at least 1"),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_metrics_tag_cardinality_limit_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						value_limit = 50
						exclude_tags = ["` + strings.Repeat("x", 101) + `"]
					}`,
				ExpectError: regexp.MustCompile("Attribute exclude_tags.* string length must be at most 100"),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_metrics_tag_cardinality_limit_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						value_limit = 50
						exclude_tags = ["x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x"]
					}`,
				ExpectError: regexp.MustCompile("Attribute exclude_tags list must contain at most 10 elements"),
			},

			// Error: action is an enum
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_metrics_tag_cardinality_limit_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						value_limit = 50
						action = "nope"
					}`,
				ExpectError: regexp.MustCompile(`Attribute action value must be one of: \["drop_tag" "drop_event"\]`),
			},

			// Error: mode is an enum
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_metrics_tag_cardinality_limit_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						value_limit = 50
						mode = "nope"
					}`,
				ExpectError: regexp.MustCompile(`Attribute mode value must be one of: \["exact" "probabilistic"\]`),
			},

			// Create with defaults
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_metrics_tag_cardinality_limit_processor" "my_processor" {
						title = "title"
						description = "desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						value_limit = 50
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_metrics_tag_cardinality_limit_processor.my_processor", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_metrics_tag_cardinality_limit_processor.my_processor", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "title",
						"description":   "desc",
						"generation_id": "0",
						"inputs.#":      "0",
						"mode":          "exact",
						"action":        "drop_event",
						"value_limit":   "50",
						"%":             "11",
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_metrics_tag_cardinality_limit_processor" "import_target" {
						title = "title"
						description = "desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						value_limit = 50
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_metrics_tag_cardinality_limit_processor.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_metrics_tag_cardinality_limit_processor.my_processor"),
				ImportStateVerify: true,
			},

			// Update fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_metrics_tag_cardinality_limit_processor" "my_processor" {
						title = "new title"
						description = "new desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						tags = ["include_me"]
						exclude_tags = ["exclude_me"]
						mode = "probabilistic"
						action = "drop_tag"
						value_limit = 150
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_metrics_tag_cardinality_limit_processor.my_processor", map[string]any{
						"pipeline_id":    "#mezmo_pipeline.test_parent.id",
						"title":          "new title",
						"description":    "new desc",
						"generation_id":  "1",
						"inputs.#":       "1",
						"inputs.0":       "#mezmo_http_source.my_source.id",
						"mode":           "probabilistic",
						"action":         "drop_tag",
						"value_limit":    "150",
						"tags.#":         "1",
						"tags.0":         "include_me",
						"exclude_tags.#": "1",
						"exclude_tags.0": "exclude_me",
						"%":              "11",
					}),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
