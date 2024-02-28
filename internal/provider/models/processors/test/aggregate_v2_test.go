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
						interval    = 3600
						strategy    = "sum"
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
						"interval":      "3600",
						"strategy":      "sum",
						"group_by.#":    "3",
						"group_by.0":    ".name",
						"group_by.1":    ".namespace",
						"group_by.2":    ".tags",
					}),
					StateDoesNotHaveFields("mezmo_aggregate_v2_processor.my_processor", []string{"script"}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_aggregate_v2_processor" "import_target" {
						title 		= "My aggregate v2 processor"
						description = "Lets aggregate stuff"
						pipeline_id = mezmo_pipeline.test_parent.id
						method      = "tumbling"
						interval    = 3600
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
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs 		= [mezmo_http_source.my_source.id]
						method      = "tumbling"
						strategy    = "sum"
						interval    = 3600
						group_by 	= [".foo", ".bar"]
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_aggregate_v2_processor.my_processor", map[string]any{
						"group_by.#":    "2",
						"group_by.0":    ".foo",
						"group_by.1":    ".bar",
						"generation_id": "1",
					}),
					StateDoesNotHaveFields("mezmo_aggregate_v2_processor.my_processor", []string{"script"}),
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
						method      = "tumbling"
						interval    = 3600
						strategy    = "sum"
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_aggregate_v2_processor.my_processor", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "new title",
						"description":   "new desc",
						"generation_id": "2",
						"inputs.#":      "1",
						"inputs.0":      "#mezmo_http_source.my_source.id",
						"method":        "tumbling",
						"interval":      "3600",
						"strategy":      "sum",
					}),
					StateDoesNotHaveFields("mezmo_aggregate_v2_processor.my_processor", []string{"script"}),
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
						strategy	= "average"
						window_duration  = 10
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_aggregate_v2_processor.my_processor", map[string]any{
						"pipeline_id":     "#mezmo_pipeline.test_parent.id",
						"title":           "new title",
						"description":     "new desc",
						"generation_id":   "3",
						"inputs.#":        "1",
						"inputs.0":        "#mezmo_http_source.my_source.id",
						"method":          "sliding",
						"strategy":        "average",
						"window_duration": "10",
					}),
					StateDoesNotHaveFields("mezmo_aggregate_v2_processor.my_processor", []string{"script"}),
				),
			},

			// Update fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_aggregate_v2_processor" "my_processor" {
						title = "new title"
						description = "new desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs 		= [mezmo_http_source.my_source.id]
						method      = "sliding"
						strategy	= "average"
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
						"generation_id":                          "4",
						"inputs.#":                               "1",
						"inputs.0":                               "#mezmo_http_source.my_source.id",
						"method":                                 "sliding",
						"strategy":                               "average",
						"window_duration":                        "10",
						"conditional.expressions.#":              "1",
						"conditional.expressions.0.field":        ".status",
						"conditional.expressions.0.operator":     "equal",
						"conditional.expressions.0.value_number": "200",
					}),
					StateDoesNotHaveFields("mezmo_aggregate_v2_processor.my_processor", []string{"script"}),
				),
			},

			// Update field
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_aggregate_v2_processor" "my_processor" {
						title 			 = "custom script"
						pipeline_id 	 = mezmo_pipeline.test_parent.id
						inputs			 = [mezmo_http_source.my_source.id]
						method 			 = "sliding"
						window_duration  = 10
						script 			 = "function process(event, metadata) { event.foo = \"bar\"; return event }"
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_aggregate_v2_processor.my_processor", map[string]any{
						"script": "function process(event, metadata) { event.foo = \"bar\"; return event }",
					}),
					StateDoesNotHaveFields("mezmo_aggregate_v2_processor.my_processor", []string{"strategy"}),
				),
			},

			// Update field
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_aggregate_v2_processor" "my_processor" {
						title 			 = "back to strategy"
						pipeline_id 	 = mezmo_pipeline.test_parent.id
						inputs			 = [mezmo_http_source.my_source.id]
						method 			 = "sliding"
						window_duration  = 10
						strategy		 = "average"
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_aggregate_v2_processor.my_processor", map[string]any{
						"strategy": "average",
					}),
					StateDoesNotHaveFields("mezmo_aggregate_v2_processor.my_processor", []string{"script"}),
				),
			},

			// Error: server-side validation for missing strategy
			{
				Config: GetCachedConfig(cacheKey) + `
				resource "mezmo_aggregate_v2_processor" "my_processor" {
					pipeline_id = mezmo_pipeline.test_parent.id
					inputs 		= [mezmo_http_source.my_source.id]
					method 		= "tumbling"
					interval 	= "36000"
					strategy    = "sum"
				}`,
				ExpectError: regexp.MustCompile("(?s)must.*be <=.*"),
			},

			// Error: server-side validation for invalid strategy
			{
				Config: GetCachedConfig(cacheKey) + `
				resource "mezmo_aggregate_v2_processor" "my_processor" {
					pipeline_id = mezmo_pipeline.test_parent.id
					inputs 	 	= [mezmo_http_source.my_source.id]
					method 	 	= "sliding"
					strategy 	= "asdf"
				}`,
				ExpectError: regexp.MustCompile("Attribute strategy value must be one of"),
			},

			// Error: server-side validation - invalid window duration
			{
				Config: GetCachedConfig(cacheKey) + `
				resource "mezmo_aggregate_v2_processor" "my_processor" {
					pipeline_id 	= mezmo_pipeline.test_parent.id
					inputs 			= [mezmo_http_source.my_source.id]
					method 			= "sliding"
					strategy 		= "average"
					window_duration = 305325235325
				}`,
				ExpectError: regexp.MustCompile("/window_duration/maximum"),
			},

			// Error: window duration
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_aggregate_v2_processor" "my_processor" {
						pipeline_id 	= mezmo_pipeline.test_parent.id
						inputs 			= [mezmo_http_source.my_source.id]
						strategy 		= "sum"
						method 			= "tumbling"
						interval 		= 10
						window_duration = 5
					}`,
				ExpectError: regexp.MustCompile("The field 'window_duration' can only be set if method == 'sliding'"),
			},

			// Error: window min
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_aggregate_v2_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs 		= [mezmo_http_source.my_source.id]
						strategy 	= "sum"
						method 		= "tumbling"
						interval 	= 10
						window_min 	= 5
					}`,
				ExpectError: regexp.MustCompile("The field 'window_min' can only be set if method == 'sliding'"),
			},

			// Error: interval
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_aggregate_v2_processor" "my_processor" {
						pipeline_id 	= mezmo_pipeline.test_parent.id
						inputs 			= [mezmo_http_source.my_source.id]
						strategy 		= "sum"
						method 			= "sliding"
						interval 		= 10
						window_duration = 5
					}`,
				ExpectError: regexp.MustCompile("The field 'interval' can only be set if method == 'tumbling'"),
			},

			// confirm manually deleted resources are recreated
			{
				Config: GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_aggregate_v2_processor" "test_processor" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 		= "new title"
					inputs 		= []
					method      = "tumbling"
  					interval    = 3600
					strategy    = "sum"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_aggregate_v2_processor.test_processor", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_aggregate_v2_processor.test_processor", "title", "new title"),
					// delete the resource
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_aggregate_v2_processor.test_processor",
					),
				),
				// verify resource will be re-created after refresh
				ExpectNonEmptyPlan: true,
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
