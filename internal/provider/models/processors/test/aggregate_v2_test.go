package processors

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestAccAggregateV2Processor(t *testing.T) {
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
						window_type = "tumbling"
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
						window_type = "tumbling"
						interval    = 3600
						operation   = "sum"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_aggregate_v2_processor.my_processor", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_aggregate_v2_processor.my_processor", map[string]any{
						"pipeline_id":     "#mezmo_pipeline.test_parent.id",
						"title":           "My aggregate v2 processor",
						"description":     "Lets aggregate stuff",
						"generation_id":   "0",
						"inputs.#":        "0",
						"window_type":     "tumbling",
						"interval":        "3600",
						"operation":       "sum",
						"event_timestamp": "timestamp",
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
						window_type = "tumbling"
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
						window_type = "tumbling"
						operation   = "sum"
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
						inputs      = [mezmo_http_source.my_source.id]
						window_type = "tumbling"
						interval    = 3600
						operation   = "sum"
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_aggregate_v2_processor.my_processor", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "new title",
						"description":   "new desc",
						"generation_id": "2",
						"inputs.#":      "1",
						"inputs.0":      "#mezmo_http_source.my_source.id",
						"window_type":   "tumbling",
						"interval":      "3600",
						"operation":     "sum",
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
						inputs      = [mezmo_http_source.my_source.id]
						window_type = "sliding"
						window_min  = 10
						operation	= "average"
						interval    = 10
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_aggregate_v2_processor.my_processor", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "new title",
						"description":   "new desc",
						"generation_id": "3",
						"inputs.#":      "1",
						"inputs.0":      "#mezmo_http_source.my_source.id",
						"window_type":   "sliding",
						"window_min":    "10",
						"operation":     "average",
						"interval":      "10",
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
						window_type = "sliding"
						window_min  = 10
						operation	= "average"
						interval    = 10
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
						"window_type":                            "sliding",
						"window_min":                             "10",
						"operation":                              "average",
						"interval":                               "10",
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
						window_type      = "sliding"
						window_min       = 10
						interval         = 10
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
						window_type 	 = "sliding"
						window_min       = 10
						interval         = 10
						operation		 = "average"
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_aggregate_v2_processor.my_processor", map[string]any{
						"operation": "average",
					}),
					StateDoesNotHaveFields("mezmo_aggregate_v2_processor.my_processor", []string{"script"}),
				),
			},

			// Error: server-side validation - invalid interval
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_aggregate_v2_processor" "my_processor" {
						pipeline_id 	= mezmo_pipeline.test_parent.id
						inputs 			= [mezmo_http_source.my_source.id]
						window_type 	= "sliding"
						window_min      = 10
						operation 		= "average"
						interval        = 305325235325
					}`,
				ExpectError: regexp.MustCompile("(?s).*/user_config/window/interval"),
			},

			// Check backend values
			{
				Config: GetCachedConfig(cacheKey) + `
						resource "mezmo_aggregate_v2_processor" "my_processor" {
							pipeline_id = mezmo_pipeline.test_parent.id
							inputs 		= [mezmo_http_source.my_source.id]
							operation 	= "sum"
							window_type = "sliding"
							window_min 	= 10
							interval 	= 10
						}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mezmo_aggregate_v2_processor.my_processor", "window_type", "sliding"),
					resource.TestCheckResourceAttr("mezmo_aggregate_v2_processor.my_processor", "window_min", "10"),
					testAccBackend("mezmo_aggregate_v2_processor.my_processor", map[string]any{
						"window": map[string]any{
							"type":       string("sliding"),
							"interval":   float64(10),
							"window_min": float64(10),
						},
						"evaluate":        map[string]any{"operation": string("SUM")},
						"event_timestamp": string("timestamp"),
					}),
				),
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
						window_type = "tumbling"
						interval    = 3600
						operation   = "sum"
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

			// Setting optional event_timestamp field
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent3" {
						title = "event timestamp test"
					}
					resource "mezmo_aggregate_v2_processor" "event_timestamp_processor" {
						pipeline_id     = mezmo_pipeline.test_parent3.id
						title           = "new processor"
						inputs          = []
						window_type     = "tumbling"
						interval        = 3600
						operation       = "sum"
						event_timestamp = ".my_ts_field"
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					StateHasExpectedValues("mezmo_aggregate_v2_processor.event_timestamp_processor", map[string]any{
						"event_timestamp": ".my_ts_field",
					}),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
