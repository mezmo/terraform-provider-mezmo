package processors

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestAggregateV2Processor(t *testing.T) {
	const cacheKey = "aggregate_v2_resources"
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
					resource "mezmo_aggregate_v2_processor" "my_processor" {
						method      = "tumbling"
  						interval    = 36000
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},

			// Create with defaults
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_aggregate_v2_processor" "my_processor" {
						title = "My aggregate v2 processor"
						description = "Lets aggregate stuff"
						pipeline_id = mezmo_pipeline.test_parent.id
						method      = "tumbling"
  						interval    = 36000
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_aggregate_v2_processor.my_processor", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_aggregate_v2_processor.my_processor", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "My aggregate v2 processor",
						"description":   "Lets aggregate stuff",
						"generation_id": "0",
						"inputs.#":      "0",
						"method":        "tumbling",
						"interval":      "36000",
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_aggregate_v2_processor" "import_target" {
						title = "My aggregate v2 processor"
						description = "Lets aggregate stuff"
						pipeline_id = mezmo_pipeline.test_parent.id
						method      = "tumbling"
  						interval    = 36000
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_aggregate_v2_processor.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_aggregate_v2_processor.my_processor"),
				ImportStateVerify: true,
			},

			// Update fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_aggregate_v2_processor" "my_processor" {
						title = "new title"
						description = "new desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						method      = "tumbling"
  						interval    = 3600
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_aggregate_v2_processor.my_processor", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "new title",
						"description":   "new desc",
						"generation_id": "1",
						"inputs.#":      "1",
						"inputs.0":      "#mezmo_http_source.my_source.id",
						"method":        "tumbling",
						"interval":      "3600",
					}),
				),
			},

			// Update fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_aggregate_v2_processor" "my_processor" {
						title = "new title"
						description = "new desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						method      = "sliding"
						strategy	= "AVG"
  						window_duration  = 10
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_aggregate_v2_processor.my_processor", map[string]any{
						"pipeline_id":     "#mezmo_pipeline.test_parent.id",
						"title":           "new title",
						"description":     "new desc",
						"generation_id":   "2",
						"inputs.#":        "1",
						"inputs.0":        "#mezmo_http_source.my_source.id",
						"method":          "sliding",
						"strategy":        "AVG",
						"window_duration": "10",
					}),
				),
			},

			// Update fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_aggregate_v2_processor" "my_processor" {
						title = "new title"
						description = "new desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						method      = "sliding"
						strategy	= "AVG"
  						window_duration  = 10
						conditional = {
							expressions = [
								{
									field = ".status"
									operator = "equal"
									value_number = 200
								}
							]
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_aggregate_v2_processor.my_processor", map[string]any{
						"pipeline_id":                            "#mezmo_pipeline.test_parent.id",
						"title":                                  "new title",
						"description":                            "new desc",
						"generation_id":                          "3",
						"inputs.#":                               "1",
						"inputs.0":                               "#mezmo_http_source.my_source.id",
						"method":                                 "sliding",
						"strategy":                               "AVG",
						"window_duration":                        "10",
						"conditional.expressions.#":              "1",
						"conditional.expressions.0.field":        ".status",
						"conditional.expressions.0.operator":     "equal",
						"conditional.expressions.0.value_number": "200",
					}),
				),
			},

			// Error: server-side validation for missing strategy
			{
				Config: GetCachedConfig(cacheKey) + `
				resource "mezmo_aggregate_v2_processor" "my_processor" {
					pipeline_id = mezmo_pipeline.test_parent.id
					inputs = [mezmo_http_source.my_source.id]
					method = "sliding"
					interval = "3600"
				}`,
				ExpectError: regexp.MustCompile("(?s)have required property.*'strategy'"),
			},

			// Error: server-side validation for invalid strategy
			{
				Config: GetCachedConfig(cacheKey) + `
				resource "mezmo_aggregate_v2_processor" "my_processor" {
					pipeline_id = mezmo_pipeline.test_parent.id
					inputs = [mezmo_http_source.my_source.id]
					method = "sliding"
					strategy = "asdf"
				}`,
				ExpectError: regexp.MustCompile("Bad Request"),
			},

			// Error: server-side validation - invalid window duration
			{
				Config: GetCachedConfig(cacheKey) + `
				resource "mezmo_aggregate_v2_processor" "my_processor" {
					pipeline_id = mezmo_pipeline.test_parent.id
					inputs = [mezmo_http_source.my_source.id]
					method = "sliding"
					strategy = "AVG"
					window_duration = 305325235325
				}`,
				ExpectError: regexp.MustCompile("/window_duration/maximum"),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
