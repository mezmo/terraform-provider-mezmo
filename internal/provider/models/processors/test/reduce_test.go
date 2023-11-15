package processors

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestReduceProcessor(t *testing.T) {
	const cacheKey = "reduce_resources"
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
					resource "mezmo_reduce_processor" "my_processor" {
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},

			// Error: expressions_group and expressions are mutually exclusive
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_reduce_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id

						flush_condition = {
							when = "ends_when"
							conditional = {
								expressions = [
									{
										field = ".status"
										operator = "equal"
										value_number = 400
									}
								],
								expressions_group = [
									{
										expressions = [
											{
												field = ".nope"
												operator = "equal"
												value_string = "nope"
											}
										],
										logical_operator = "AND"
									}
								]
							}
						}
					}`,
				ExpectError: regexp.MustCompile("(?s)Attribute \"flush_condition.conditional.expressions_group\" cannot be specified.*when \"flush_condition.conditional.expressions\" is specified"),
			},

			// Error: value_string and value_number are mutually exclusive
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_reduce_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id

						flush_condition = {
							when = "ends_when"
							conditional = {
								expressions = [
									{
										field = ".status"
										operator = "equal"
										value_number = 400
										value_string = "NOPE"
									}
								]
							}
						}
					}`,
				ExpectError: regexp.MustCompile("(?s)Attribute.*value_number.*cannot be.*value_string"),
			},

			// Defaults only
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_reduce_processor" "default_values" {
						title = "processor title"
						description = "processor desc"
						pipeline_id = mezmo_pipeline.test_parent.id
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_reduce_processor.default_values", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_reduce_processor.default_values", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "processor title",
						"description":   "processor desc",
						"generation_id": "0",
						"inputs.#":      "0",
						"duration_ms":   "30000",
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_reduce_processor" "import_target" {
						title = "processor title"
						description = "processor desc"
						pipeline_id = mezmo_pipeline.test_parent.id
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_reduce_processor.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_reduce_processor.default_values"),
				ImportStateVerify: true,
			},

			// Add options
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_reduce_processor" "default_values" {
						title = "processor title"
						description = "processor desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]

						duration_ms = 15000
						group_by    = [".field1.id"]
						date_formats = [
							{
								field = ".datetime"
								format = "%d/%m/%Y:%T"
							}
						]
						merge_strategies = [
							{
								field    = ".method"
								strategy = "array"
							}
						]
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_reduce_processor.default_values", map[string]any{
						"pipeline_id":                 "#mezmo_pipeline.test_parent.id",
						"title":                       "processor title",
						"description":                 "processor desc",
						"generation_id":               "1",
						"inputs.#":                    "1",
						"duration_ms":                 "15000",
						"group_by.#":                  "1",
						"group_by.0":                  ".field1.id",
						"date_formats.#":              "1",
						"date_formats.0.field":        ".datetime",
						"date_formats.0.format":       "%d/%m/%Y:%T",
						"merge_strategies.#":          "1",
						"merge_strategies.0.field":    ".method",
						"merge_strategies.0.strategy": "array",
					}),
				),
			},

			// Single level flush_condition, default logical_operation
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_reduce_processor" "my_processor" {
						title = "processor title"
						description = "processor desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]

						flush_condition = {
							when = "starts_when"
							conditional = {
								expressions = [
									{
										field = ".status"
										operator = "equal"
										value_number = 400
									}
								]
							},
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_reduce_processor.my_processor", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_reduce_processor.my_processor", map[string]any{
						"pipeline_id":          "#mezmo_pipeline.test_parent.id",
						"title":                "processor title",
						"description":          "processor desc",
						"generation_id":        "0",
						"inputs.#":             "1",
						"inputs.0":             "#mezmo_http_source.my_source.id",
						"flush_condition.when": "starts_when",
						"flush_condition.conditional.expressions.#":              "1",
						"flush_condition.conditional.expressions.0.field":        ".status",
						"flush_condition.conditional.expressions.0.operator":     "equal",
						"flush_condition.conditional.expressions.0.value_number": "400",
						"flush_condition.conditional.logical_operation":          "AND",
					}),
				),
			},

			// Single level flush_condition, provides logical_operation
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_reduce_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]

						flush_condition = {
							when = "starts_when"
							conditional = {
								expressions = [
									{
										field = ".level"
										operator = "equal"
										value_string = "ERROR"
									},
									{
										field = ".app"
										operator = "equal"
										value_string = "worker"
									}
								],
								logical_operation = "OR"
							},
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_reduce_processor.my_processor", map[string]any{
						"generation_id": "1",
						"flush_condition.conditional.expressions.#":              "2",
						"flush_condition.conditional.expressions.0.field":        ".level",
						"flush_condition.conditional.expressions.0.operator":     "equal",
						"flush_condition.conditional.expressions.0.value_string": "ERROR",
						"flush_condition.conditional.expressions.1.field":        ".app",
						"flush_condition.conditional.expressions.1.operator":     "equal",
						"flush_condition.conditional.expressions.1.value_string": "worker",
						"flush_condition.conditional.logical_operation":          "OR",
					}),
				),
			},

			// Nested flush_condition conditionals
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_reduce_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]

						flush_condition = {
							when = "starts_when"
							conditional = {
								expressions_group = [
									{
										expressions = [
											{
												field = ".level"
												operator = "equal"
												value_string = "ERROR"
											},
											{
												field = ".app"
												operator = "equal"
												value_string = "worker"
											}
										],
										logical_operation = "AND"
									},
									{
										expressions = [
											{
												field = ".status"
												operator = "equal"
												value_number = 500
											},
											{
												field = ".app"
												operator = "equal"
												value_string = "service"
											}
										],
										logical_operation = "AND"
									},
									{
										expressions_group = [
											{
												expressions = [
													{
														field = ".deeper"
														operator = "equal"
														value_string = "yep"
													},
													{
														field = ".other.deeper"
														operator = "equal"
														value_string = "getting deep now"
													}
												],
												logical_operation = "AND"
											}
										],
										logical_operation = "AND"
									},
								]
								logical_operation = "OR"
							}
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_reduce_processor.my_processor", map[string]any{
						"generation_id": "2",
						"flush_condition.conditional.expressions_group.#":                                                "3",
						"flush_condition.conditional.expressions_group.0.expressions.#":                                  "2",
						"flush_condition.conditional.expressions_group.0.expressions.0.field":                            ".level",
						"flush_condition.conditional.expressions_group.0.expressions.0.operator":                         "equal",
						"flush_condition.conditional.expressions_group.0.expressions.0.value_string":                     "ERROR",
						"flush_condition.conditional.expressions_group.0.expressions.1.field":                            ".app",
						"flush_condition.conditional.expressions_group.0.expressions.1.operator":                         "equal",
						"flush_condition.conditional.expressions_group.0.expressions.1.value_string":                     "worker",
						"flush_condition.conditional.expressions_group.0.logical_operation":                              "AND",
						"flush_condition.conditional.expressions_group.1.expressions.0.field":                            ".status",
						"flush_condition.conditional.expressions_group.1.expressions.0.operator":                         "equal",
						"flush_condition.conditional.expressions_group.1.expressions.0.value_number":                     "500",
						"flush_condition.conditional.expressions_group.1.expressions.1.field":                            ".app",
						"flush_condition.conditional.expressions_group.1.expressions.1.operator":                         "equal",
						"flush_condition.conditional.expressions_group.1.expressions.1.value_string":                     "service",
						"flush_condition.conditional.expressions_group.1.logical_operation":                              "AND",
						"flush_condition.conditional.expressions_group.2.expressions_group.#":                            "1",
						"flush_condition.conditional.expressions_group.2.expressions_group.0.expressions.#":              "2",
						"flush_condition.conditional.expressions_group.2.expressions_group.0.expressions.0.field":        ".deeper",
						"flush_condition.conditional.expressions_group.2.expressions_group.0.expressions.0.operator":     "equal",
						"flush_condition.conditional.expressions_group.2.expressions_group.0.expressions.0.value_string": "yep",
						"flush_condition.conditional.expressions_group.2.expressions_group.0.expressions.1.field":        ".other.deeper",
						"flush_condition.conditional.expressions_group.2.expressions_group.0.expressions.1.operator":     "equal",
						"flush_condition.conditional.expressions_group.2.expressions_group.0.expressions.1.value_string": "getting deep now",
						"flush_condition.conditional.expressions_group.2.expressions_group.0.logical_operation":          "AND",
						"flush_condition.conditional.expressions_group.2.logical_operation":                              "AND",
						"flush_condition.conditional.logical_operation":                                                  "OR",
					}),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
