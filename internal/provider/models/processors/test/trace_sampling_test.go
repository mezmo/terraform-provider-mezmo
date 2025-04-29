package processors

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/providertest"
)

func TestAccTraceSamplingErrors(t *testing.T) {
	const cacheKey = "tracesampling_errors_cache_key"

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
					resource "mezmo_trace_sampling_processor" "my_processor" {
						sample_type = "head"
						trace_id_field = ".trace_id"
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},
			// Error: sample_type is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_trace_sampling_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile("The argument \"sample_type\" is required"),
			},
		},
	})
}

func TestAccTraceHeadSamplingSuccess(t *testing.T) {
	const cacheKey = "trace_head_sampling_success_cache_key"

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
					}`),
			},
			// parent span isn't applicable
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_trace_sampling_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						sample_type = "head"
						trace_id_field = ".trace_id"
						rate = 10
						parent_span_id_field = ".parent_span_id"
					}`,
				ExpectError: regexp.MustCompile("Attribute \"parent_span_id_field\" is not applicable for head sampling."),
			},
			// conditionals span isn't applicable
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_trace_sampling_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						sample_type = "head"
						trace_id_field = ".trace_id"
						rate = 10
						conditionals = [{
							rate = 10
							conditional = {}
						}]
					}`,
				ExpectError: regexp.MustCompile("Attribute \"conditionals\" is not applicable for head sampling."),
			},
			// rate must be >= 2
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_trace_sampling_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						sample_type = "head"
						trace_id_field = ".trace_id"
						rate = 1
					}`,
				ExpectError: regexp.MustCompile("(?s)Attribute rate value must be at least 2"),
			},
			// successful head based sampling with defaults
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_trace_sampling_processor" "head_based_sampler" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "processor title"
						description = "processor desc"
						inputs = [mezmo_http_source.my_source.id]

						sample_type = "head"
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_trace_sampling_processor.head_based_sampler", map[string]any{
						"pipeline_id":          "#mezmo_pipeline.test_parent.id",
						"title":                "processor title",
						"description":          "processor desc",
						"generation_id":        "0",
						"inputs.#":             "1",
						"sample_type":          "head",
						"rate":                 "10",
						"trace_id_field":       ".trace_id",
						"parent_span_id_field": nil, // used in tail sampling
						"conditionals":         nil, // used in tail sampling
					}),
				)},
			// update
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_trace_sampling_processor" "head_based_sampler" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "processor title"
						description = "processor desc"
						inputs = [mezmo_http_source.my_source.id]

						sample_type = "head"
						rate = 15
						trace_id_field = ".trace_id_y"
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_trace_sampling_processor.head_based_sampler", map[string]any{
						"pipeline_id":          "#mezmo_pipeline.test_parent.id",
						"title":                "processor title",
						"description":          "processor desc",
						"generation_id":        "1",
						"inputs.#":             "1",
						"sample_type":          "head",
						"rate":                 "15",
						"trace_id_field":       ".trace_id_y",
						"parent_span_id_field": nil, // used in tail sampling
						"conditionals":         nil, // used in tail sampling
					}),
				)},
		},
	})
}

func TestAccTraceTailSamplingSuccess(t *testing.T) {
	const cacheKey = "trace_tail_sampling_success_cache_key"

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
					}`),
			},
			// rate isn't applicable
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_trace_sampling_processor" "my_processor_rate" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "processor title"
						description = "processor desc"
						inputs = [mezmo_http_source.my_source.id]

						sample_type = "tail"
						rate = 10
						conditionals = [
						{
							rate = 15
							conditional = {
								expressions = [
									{
										field = ".status"
										operator = "equal"
										value_number = 200
										negate = false
									}
								]
								logical_operation = "AND"
							},
							_output_name = "ab1869cd"
						},
					]
				}`,
				ExpectError: regexp.MustCompile("Attribute \"rate\" is not applicable for tail sampling."),
			},
			// missing conditionals
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_trace_sampling_processor" "my_processor_conditional" {
						pipeline_id = mezmo_pipeline.test_parent.id
						sample_type = "tail"
					}`,
				ExpectError: regexp.MustCompile("Attribute \"conditionals\" is required for tail sampling."),
			},
			// successful tail based sampling with defaults
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_trace_sampling_processor" "tail_based_sampler" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "processor title"
						description = "processor desc"
						inputs = [mezmo_http_source.my_source.id]

						sample_type = "tail"
						conditionals = [
						{
							rate = 15
							conditional = {
								expressions = [
									{
										field = ".status"
										operator = "equal"
										value_number = 200
										negate = false
									}
								]
								logical_operation = "AND"
							},
							_output_name = "ab1869cd"
						},
					]
				}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_trace_sampling_processor.tail_based_sampler", map[string]any{
						"pipeline_id":                 "#mezmo_pipeline.test_parent.id",
						"title":                       "processor title",
						"description":                 "processor desc",
						"generation_id":               "0",
						"inputs.#":                    "1",
						"sample_type":                 "tail",
						"rate":                        nil, // used in head sampling
						"trace_id_field":              ".trace_id",
						"parent_span_id_field":        ".parent_span_id",
						"conditionals.#":              "1",
						"conditionals.0.rate":         "15",
						"conditionals.0._output_name": "ab1869cd",
						"conditionals.0.conditional.expressions.#":              "1",
						"conditionals.0.conditional.expressions.0.field":        ".status",
						"conditionals.0.conditional.expressions.0.operator":     "equal",
						"conditionals.0.conditional.expressions.0.value_number": "200",
						"conditionals.0.conditional.expressions.0.value_string": nil,
						"conditionals.0.conditional.expressions.0.negate":       "false",
						"conditionals.0.conditional.expressions_group":          nil,
						"conditionals.0.conditional.logical_operation":          "AND",
					}),
				)},
			// update
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_trace_sampling_processor" "tail_based_sampler" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "processor title"
						description = "processor desc"
						inputs = [mezmo_http_source.my_source.id]

						sample_type = "tail"
						trace_id_field = ".trace_id_z"
						parent_span_id_field = ".parent_span_id_z"
						conditionals = [
						{
							rate = 20
							conditional = {
								expressions = [
									{
										field = ".status_code"
										operator = "greater_or_equal"
										value_number = 500
										negate = true
									}
								]
								logical_operation = "AND"
							},
							_output_name = "ab1869cd"
						}
					]
				}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_trace_sampling_processor.tail_based_sampler", map[string]any{
						"pipeline_id":                 "#mezmo_pipeline.test_parent.id",
						"title":                       "processor title",
						"description":                 "processor desc",
						"generation_id":               "1",
						"inputs.#":                    "1",
						"sample_type":                 "tail",
						"rate":                        nil, // used in head sampling
						"trace_id_field":              ".trace_id_z",
						"parent_span_id_field":        ".parent_span_id_z",
						"conditionals.#":              "1",
						"conditionals.0.rate":         "20",
						"conditionals.0._output_name": "ab1869cd",
						"conditionals.0.conditional.expressions.#":              "1",
						"conditionals.0.conditional.expressions.0.field":        ".status_code",
						"conditionals.0.conditional.expressions.0.operator":     "greater_or_equal",
						"conditionals.0.conditional.expressions.0.value_number": "500",
						"conditionals.0.conditional.expressions.0.value_string": nil,
						"conditionals.0.conditional.expressions.0.negate":       "true",
						"conditionals.0.conditional.expressions_group":          nil,
						"conditionals.0.conditional.logical_operation":          "AND",
					}),
				)},
		},
	})
}
