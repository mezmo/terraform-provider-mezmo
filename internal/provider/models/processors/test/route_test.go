package processors

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestAccRouteProcessor(t *testing.T) {
	const cacheKey = "route_resources"
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
					resource "mezmo_route_processor" "my_processor" {
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},

			// Error: conditionals is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_route_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile("The argument \"conditionals\" is required"),
			},

			// Error: expressions_group and expressions are mutually exclusive
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_route_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id

						conditionals = [
							{
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
								label = "Welcome messages"
							}
						]
					}`,
				ExpectError: regexp.MustCompile("(?s)Attribute \"conditionals\\[0\\].expressions_group\" cannot be specified when.*\"conditionals\\[0\\].expressions\" is specified"),
			},

			// Error: label is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_route_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id

						conditionals = [
							{
								expressions = [
									{
										field = ".status"
										operator = "equal"
										value_number = 200
									}
								]
								label = "Welcome messages"
							},
							{
								expressions = [
									{
										field = ".level"
										operator = "not_equal"
										value_string = "error"
									}
								]
							}
						]
					}`,
				ExpectError: regexp.MustCompile("(?s)Inappropriate value for attribute \"conditionals\": element 1: attribute.*\"label\" is required"),
			},

			// Error: label length too long
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_route_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id

						conditionals = [
							{
								expressions = [
									{
										field = ".status"
										operator = "equal"
										value_number = 200
									}
								]
								label = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
							}
						]
					}`,
				ExpectError: regexp.MustCompile("(?s).*label string length must be at most 255.*"),
			},

			// Error: value_string and value_number are mutually exclusive
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_route_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id

						conditionals = [
							{
								expressions = [
									{
										field = ".status"
										operator = "equal"
										value_number = 400
										value_string = "NOPE"
									}
								]
								label = "invalid"
							}
						]
					}`,
				ExpectError: regexp.MustCompile("(?s)Attribute.*value_number.*cannot be.*value_string"),
			},

			// Single expression
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_route_processor" "single_expression" {
						title = "processor title"
						description = "processor desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]

						conditionals = [
							{
								expressions = [
									{
										field = ".status"
										operator = "equal"
										value_number = 200
									}
								]
								label = "success logs"
							}
						]
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_route_processor.single_expression", map[string]any{
						"pipeline_id":                               "#mezmo_pipeline.test_parent.id",
						"title":                                     "processor title",
						"description":                               "processor desc",
						"generation_id":                             "0",
						"inputs.#":                                  "1",
						"conditionals.#":                            "1",
						"conditionals.0.label":                      "success logs",
						"conditionals.0.expressions.#":              "1",
						"conditionals.0.expressions.0.field":        ".status",
						"conditionals.0.expressions.0.operator":     "equal",
						"conditionals.0.expressions.0.value_number": "200",
						"conditionals.0.output_name":                regexp.MustCompile(`^.+\..+$`),
						"unmatched":                                 regexp.MustCompile(`^.+\..+$`),
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_route_processor" "import_target" {
						title = "processor title"
						description = "processor desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]

						conditionals = [
							{
								expressions = [
									{
										field = ".status"
										operator = "equal"
										value_number = 200
									}
								]
								label = "success logs"
							}
						]
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_route_processor.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_route_processor.single_expression"),
				ImportStateVerify: true,
			},

			// Nested expression
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_route_processor" "nested_expression" {
						title = "processor title"
						description = "processor desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]

						conditionals = [
							{
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
								label = "error logs"
							}
						]
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_route_processor.nested_expression", map[string]any{
						"pipeline_id":                        "#mezmo_pipeline.test_parent.id",
						"title":                              "processor title",
						"description":                        "processor desc",
						"generation_id":                      "0",
						"inputs.#":                           "1",
						"conditionals.#":                     "1",
						"conditionals.0.label":               "error logs",
						"conditionals.0.expressions_group.#": "2",
						"conditionals.0.expressions_group.0.expressions.#":                                  "2",
						"conditionals.0.expressions_group.0.expressions.0.field":                            ".label",
						"conditionals.0.expressions_group.0.expressions.0.operator":                         "equal",
						"conditionals.0.expressions_group.0.expressions.0.value_string":                     "account",
						"conditionals.0.expressions_group.0.expressions.1.field":                            ".app_name",
						"conditionals.0.expressions_group.0.expressions.1.operator":                         "ends_with",
						"conditionals.0.expressions_group.0.expressions.1.value_string":                     "service",
						"conditionals.0.expressions_group.1.expressions_group.#":                            "1",
						"conditionals.0.expressions_group.1.expressions_group.0.expressions.#":              "2",
						"conditionals.0.expressions_group.1.expressions_group.0.expressions.0.field":        ".level",
						"conditionals.0.expressions_group.1.expressions_group.0.expressions.0.operator":     "greater_or_equal",
						"conditionals.0.expressions_group.1.expressions_group.0.expressions.0.value_number": "300",
						"conditionals.0.expressions_group.1.expressions_group.0.expressions.1.field":        ".tag",
						"conditionals.0.expressions_group.1.expressions_group.0.expressions.1.operator":     "contains",
						"conditionals.0.expressions_group.1.expressions_group.0.expressions.1.value_string": "error",
						"conditionals.0.output_name":                                                        regexp.MustCompile(`^.+\..+$`),
						"unmatched":                                                                         regexp.MustCompile(`^.+\..+$`),
					}),
				),
			},

			// Multiple conditionals
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_route_processor" "multiple_conditionals" {
						title = "processor title"
						description = "processor desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]

						conditionals = [
							{
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
									}
								]
								label = "app logs"
							},
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
								label = "error logs"
							}
						]
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_route_processor.multiple_conditionals", map[string]any{
						"pipeline_id":                        "#mezmo_pipeline.test_parent.id",
						"title":                              "processor title",
						"description":                        "processor desc",
						"generation_id":                      "0",
						"inputs.#":                           "1",
						"conditionals.#":                     "2",
						"conditionals.0.label":               "app logs",
						"conditionals.1.label":               "error logs",
						"conditionals.0.expressions_group.#": "1",
						"conditionals.0.expressions_group.0.expressions.#":              "2",
						"conditionals.0.expressions_group.0.expressions.0.field":        ".label",
						"conditionals.0.expressions_group.0.expressions.0.operator":     "equal",
						"conditionals.0.expressions_group.0.expressions.0.value_string": "account",
						"conditionals.0.expressions_group.0.expressions.1.field":        ".app_name",
						"conditionals.0.expressions_group.0.expressions.1.operator":     "ends_with",
						"conditionals.0.expressions_group.0.expressions.1.value_string": "service",
						"conditionals.0.output_name":                                    regexp.MustCompile(`^.+\..+$`),
						"conditionals.1.expressions.#":                                  "2",
						"conditionals.1.expressions.0.field":                            ".level",
						"conditionals.1.expressions.0.operator":                         "greater_or_equal",
						"conditionals.1.expressions.0.value_number":                     "300",
						"conditionals.1.expressions.1.field":                            ".tag",
						"conditionals.1.expressions.1.operator":                         "contains",
						"conditionals.1.expressions.1.value_string":                     "error",
						"conditionals.1.output_name":                                    regexp.MustCompile(`^.+\..+$`),
						"unmatched":                                                     regexp.MustCompile(`^.+\..+$`),
					}),
				),
			},

			// Connect outputs to destinations
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_route_processor" "with_outputs" {
						title = "processor title"
						description = "processor desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]

						conditionals = [
							{
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
								label = "service logs"
							},
							{
								expressions = [
									{
										field = ".level"
										operator = "equal"
										value_string = "error"
									}
								]
								label = "error logs"
							}
						]
					}
					resource "mezmo_blackhole_destination" "destination1" {
						title = "blackhole title"
						description = "blackhole desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_route_processor.with_outputs.conditionals.0.output_name]
					}
					resource "mezmo_logs_destination" "destination2" {
						title = "logs title"
						description = "logs desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_route_processor.with_outputs.conditionals.1.output_name]
						ingestion_key = "my_key"
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_blackhole_destination.destination1", map[string]any{
						"pipeline_id": "#mezmo_pipeline.test_parent.id",
						"inputs.#":    "1",
					}),
					StateHasExpectedValues("mezmo_logs_destination.destination2", map[string]any{
						"pipeline_id": "#mezmo_pipeline.test_parent.id",
						"inputs.#":    "1",
					}),
				),
			},
			// confirm manually deleted resources are recreated
			{
				Config: GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_route_processor" "test_processor" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					inputs 			= []
					conditionals = [
						{
							expressions = [
								{
									field = ".status"
									operator = "equal"
									value_number = 200
								}
							]
							label = "success logs"
						}
					]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_route_processor.test_processor", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_route_processor.test_processor", "title", "new title"),
					// delete the resource
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_route_processor.test_processor",
					),
				),
				// verify resource will be re-created after refresh
				ExpectNonEmptyPlan: true,
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
