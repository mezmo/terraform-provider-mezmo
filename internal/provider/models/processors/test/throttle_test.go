package processors

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/v5/internal/provider/providertest"
)

func TestAccThrottleProcessor(t *testing.T) {
	const cacheKey = "throttle_resourcess"
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
					resource "mezmo_throttle_processor" "my_processor" {
						threshold = 100
						window_ms = 10000
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},

			// Error: `threshold` is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_throttle_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						window_ms = 10000
					}`,
				ExpectError: regexp.MustCompile("The argument \"threshold\" is required"),
			},

			// Error: `window_ms` is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_throttle_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						threshold = 100
					}`,
				ExpectError: regexp.MustCompile("The argument \"window_ms\" is required"),
			},

			// Error: `window_ms` values validates length
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_throttle_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						threshold = 100
						window_ms = 10
					}`,
				ExpectError: regexp.MustCompile("Attribute window_ms value must be at least 1000"),
			},

			// Create with default required params
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_throttle_processor" "my_processor" {
						title = "My throttle processor"
						description = "Lets throttle stuff"
						pipeline_id = mezmo_pipeline.test_parent.id
						threshold = 100
						window_ms = 1000
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_throttle_processor.my_processor", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_throttle_processor.my_processor", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "My throttle processor",
						"description":   "Lets throttle stuff",
						"generation_id": "0",
						"inputs.#":      "0",
						"threshold":     "100",
						"window_ms":     "1000",
					}),
				),
			},

			// Create with defaults and optionals
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_throttle_processor" "my_processor" {
						title = "My throttle processor"
						description = "Lets throttle stuff"
						inputs = [mezmo_http_source.my_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id
						threshold = 100
						window_ms = 1000
						key_field = ".my_field"
						exclude = {
							expressions = [
								{
									field = ".status"
									operator = "equal"
									value_number = 500
									negate = true
								}
							]
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_throttle_processor.my_processor", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_throttle_processor.my_processor", map[string]any{
						"pipeline_id":                        "#mezmo_pipeline.test_parent.id",
						"title":                              "My throttle processor",
						"description":                        "Lets throttle stuff",
						"generation_id":                      "1",
						"inputs.#":                           "1",
						"threshold":                          "100",
						"window_ms":                          "1000",
						"key_field":                          ".my_field",
						"exclude.expressions.#":              "1",
						"exclude.expressions.0.field":        ".status",
						"exclude.expressions.0.operator":     "equal",
						"exclude.expressions.0.value_number": "500",
						"exclude.expressions.0.negate":       "true",
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_throttle_processor" "import_target" {
						title = "My throttle processor"
						description = "Lets throttle stuff"
						pipeline_id = mezmo_pipeline.test_parent.id
						threshold = 100
						window_ms = 1000
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_throttle_processor.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_throttle_processor.my_processor"),
				ImportStateVerify: true,
			},

			// Update fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_throttle_processor" "my_processor" {
						title = "new title"
						description = "new desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						threshold = 200
						window_ms = 2000
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_throttle_processor.my_processor", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "new title",
						"description":   "new desc",
						"generation_id": "2",
						"inputs.#":      "1",
						"inputs.0":      "#mezmo_http_source.my_source.id",
						"threshold":     "200",
						"window_ms":     "2000",
					}),
				),
			},

			//Error: server-side validation - key field invalid syntax
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_throttle_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						threshold = 10
						window_ms = 1000
						key_field = "10"
					}`,
				ExpectError: regexp.MustCompile("(?s).*/user_config/key_field"),
			},
			//Error: server-side validation - bad conditional operator
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_throttle_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						threshold = 1
						window_ms = 1000
						key_field = ".test"
						exclude = {
							expressions = [
								{
									field = ".status"
									operator = "equalto"
									value_number = 200
									negate = true
								}
							]
						}
					}`,
				ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
			},

			// Check backend values
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_throttle_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						threshold = 100
						window_ms = 1000
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mezmo_throttle_processor.my_processor", "threshold", "100"),
					resource.TestCheckResourceAttr("mezmo_throttle_processor.my_processor", "window_ms", "1000"),
					testAccBackend("mezmo_throttle_processor.my_processor", map[string]any{
						"threshold": float64(100),
						"window_ms": float64(1000),
					}),
				),
			},

			// confirm manually deleted resources are recreated
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent2" {
						title = "pipeline"
					}
					resource "mezmo_throttle_processor" "test_processor" {
						pipeline_id = mezmo_pipeline.test_parent2.id
						title = "new title"
						inputs = []
						threshold = 100
						window_ms = 1000
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_throttle_processor.test_processor", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_throttle_processor.test_processor", "title", "new title"),
					// delete the resource
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_throttle_processor.test_processor",
					),
				),
				// verify resource will be re-created after refresh
				ExpectNonEmptyPlan: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
