package processors

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/providertest"
)

func TestAccFilterProcessor(t *testing.T) {
	const cacheKey = "filter_resources"
	SetCachedConfig(cacheKey, `
		resource "mezmo_pipeline" "test_parent" {
			title = "pipeline"
		}
		resource "mezmo_http_source" "my_source" {
			pipeline_id = mezmo_pipeline.test_parent.id
		}`,
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Error: pipeline_id is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_filter_processor" "my_processor" {
						action = "allow"
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
				ExpectError: regexp.MustCompile(`The argument "pipeline_id" is required`),
			},

			// Error: `action` is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_filter_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
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
				ExpectError: regexp.MustCompile(`The argument "action" is required, but no definition was found`),
			},
			// Error: `action` must be one of the allowed values
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_filter_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						action = "unknown"
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
				ExpectError: regexp.MustCompile(`Attribute action value must be one of: \["allow" "drop"\], got: "unknown"`),
			},
			// Error: `conditional` is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_filter_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						action = "allow"
					}`,
				ExpectError: regexp.MustCompile(`The argument "conditional" is required, but no definition was found`),
			},
			// Error: `conditional` expressions and expressions_group must be mutually exclusive
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_filter_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						action = "allow"
						conditional = {
							expressions = [
								{
									field = ".status"
									operator = "equal"
									value_number = 200
								}
							]
							expressions_group = [
								{
									expressions = [
										{
											field = ".status"
											operator = "equal"
											value_number = 200
										}
									]
									logical_operation = "AND"
								}
							]
						}
					}`,
				ExpectError: regexp.MustCompile(`(?s)Attribute "conditional.expressions_group" cannot be specified when.*"conditional.expressions" is specified`),
			},
			// Error: `value_string` and `value_number` are mutually exclusive
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_filter_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						action = "allow"
						conditional = {
							expressions = [
								{
									field = ".status"
									operator = "equal"
									value_number = 200
									value_string = "welcome"
								}
							]
						}
					}`,
				ExpectError: regexp.MustCompile(`(?s)Attribute "conditional.expressions\[0\].value_number" cannot be specified when.*"conditional.expressions\[0\].value_string" is specified`),
			},
			// Single expression
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_filter_processor" "single_expression" {
						title = "processor title"
						description = "processor desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						action = "allow"
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
					StateHasExpectedValues("mezmo_filter_processor.single_expression", map[string]any{
						"pipeline_id":                            "#mezmo_pipeline.test_parent.id",
						"title":                                  "processor title",
						"description":                            "processor desc",
						"generation_id":                          "0",
						"inputs.#":                               "1",
						"conditional.expressions.#":              "1",
						"conditional.expressions.0.field":        ".status",
						"conditional.expressions.0.operator":     "equal",
						"conditional.expressions.0.value_number": "200",
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_filter_processor" "import_target" {
				title = "processor title"
						description = "processor desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						action = "allow"
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
				ImportState:       true,
				ResourceName:      "mezmo_filter_processor.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_filter_processor.single_expression"),
				ImportStateVerify: true,
			},

			// Complex expression
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_filter_processor" "nested_expression" {
						title = "processor title"
						description = "processor desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						action = "allow"
						conditional = {
							expressions_group = [
								{
									expressions = [
										{
											field        = ".label"
											operator     = "equal"
											value_string = "account"
										},
										{
											field        = ".app_name"
											operator     = "ends_with"
											value_string = "service"
										},
									]
									logical_operation = "OR"
								},
								{
									expressions_group = [
										{
											expressions = [
												{
													field = ".level"
													operator = "greater_or_equal"
													value_number = 300
												},
												{
													field = ".tag"
													operator = "contains"
													value_string = "error"
												}
											]
											logical_operation = "OR"
										}
									]
								}
							]
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_filter_processor.nested_expression", map[string]any{
						"pipeline_id":                     "#mezmo_pipeline.test_parent.id",
						"title":                           "processor title",
						"description":                     "processor desc",
						"generation_id":                   "0",
						"inputs.#":                        "1",
						"conditional.expressions_group.#": "2",
						"conditional.expressions_group.0.expressions.#":                                  "2",
						"conditional.expressions_group.0.logical_operation":                              "OR",
						"conditional.expressions_group.0.expressions.0.field":                            ".label",
						"conditional.expressions_group.0.expressions.0.operator":                         "equal",
						"conditional.expressions_group.0.expressions.0.value_string":                     "account",
						"conditional.expressions_group.0.expressions.1.field":                            ".app_name",
						"conditional.expressions_group.0.expressions.1.operator":                         "ends_with",
						"conditional.expressions_group.0.expressions.1.value_string":                     "service",
						"conditional.expressions_group.1.expressions_group.0.expressions.#":              "2",
						"conditional.expressions_group.1.expressions_group.0.logical_operation":          "OR",
						"conditional.expressions_group.1.expressions_group.0.expressions.0.field":        ".level",
						"conditional.expressions_group.1.expressions_group.0.expressions.0.operator":     "greater_or_equal",
						"conditional.expressions_group.1.expressions_group.0.expressions.0.value_number": "300",
						"conditional.expressions_group.1.expressions_group.0.expressions.1.field":        ".tag",
						"conditional.expressions_group.1.expressions_group.0.expressions.1.operator":     "contains",
						"conditional.expressions_group.1.expressions_group.0.expressions.1.value_string": "error",
					}),
				),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_filter_processor" "with_destination" {
						title = "processor title"
						description = "processor desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						action = "drop"
						conditional = {
							expressions_group = [
								{
									expressions = [
										{
											field        = ".label"
											operator     = "equal"
											value_string = "account"
										},
										{
											field        = ".app_name"
											operator     = "ends_with"
											value_string = "service"
										},
									]
								},
							]
						}
					}
					resource "mezmo_blackhole_destination" "destination1" {
						title = "blackhole title"
						description = "blackhole desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_filter_processor.with_destination.id]
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_blackhole_destination.destination1", map[string]any{
						"pipeline_id": "#mezmo_pipeline.test_parent.id",
						"inputs.#":    "1",
						"inputs.0":    "#mezmo_filter_processor.with_destination.id",
					}),
				),
			},
			// confirm manually deleted resources are recreated
			{
				Config: GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_filter_processor" "test_processor" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					inputs 			= []
					action 			= "allow"
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
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_filter_processor.test_processor", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_filter_processor.test_processor", "title", "new title"),
					// delete the resource
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_filter_processor.test_processor",
					),
				),
				// verify resource will be re-created after refresh
				ExpectNonEmptyPlan: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
