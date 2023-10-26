package processors

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestEventToMetricProcessor(t *testing.T) {
	const cacheKey = "event_to_metric_processor_resources"
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
					resource "mezmo_event_to_metric_processor" "my_processor" {
					}`,
				ExpectError: regexp.MustCompile(`The argument "pipeline_id" is required`),
			},

			// Error: required fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_event_to_metric_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile(`The argument "metric_name" is required`),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_event_to_metric_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile(`The argument "metric_kind" is required`),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_event_to_metric_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile(`The argument "metric_type" is required`),
			},

			// Error: regex violations
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_event_to_metric_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						metric_name = "^nope!"
						metric_type = "counter"
						metric_kind = "absolute"
						namespace_field = ".namespace"
					}`,
				ExpectError: regexp.MustCompile(`Attribute metric_name has invalid characters; See documention`),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_event_to_metric_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						metric_name = "my_tag"
						metric_type = "counter"
						metric_kind = "absolute"
						namespace_field = ".namespace"
						tags = [
							{
								name = "no:colon"
								value_type = "field"
								value = ".nope"
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`(?s)Attribute tags\[0\].name has invalid characters; See.*documention`),
			},

			// Error: length violations
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_event_to_metric_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						metric_name = "` + strings.Repeat("x", 129) + `"
						metric_type = "counter"
						metric_kind = "absolute"
						namespace_field = ".namespace"
					}`,
				ExpectError: regexp.MustCompile(`Attribute metric_name string length must be at most 128`),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_event_to_metric_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						metric_name = "my_metric"
						metric_type = "counter"
						metric_kind = "absolute"
						namespace_field = ".namespace"
						value_field = ""
					}`,
				ExpectError: regexp.MustCompile(`Attribute value_field string length must be at least 1`),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_event_to_metric_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						metric_name = "my_metric"
						metric_type = "counter"
						metric_kind = "absolute"
						namespace_field = ".namespace"
						value_field = ".something"
						tags = [
							{
								name = ""
								value_type = "field"
								value = ".something_else"
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`Attribute tags\[0\].name string length must be at least 1`),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_event_to_metric_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						metric_name = "my_metric"
						metric_type = "counter"
						metric_kind = "absolute"
						namespace_field = ".namespace"
						value_field = ".something"
						tags = [
							{
								name = "` + strings.Repeat("x", 129) + `"
								value_type = "field"
								value = ".something_else"
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`Attribute tags\[0\].name string length must be at most 128`),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_event_to_metric_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						metric_name = "my_metric"
						metric_type = "counter"
						metric_kind = "absolute"
						namespace_field = ".namespace"
						value_field = ".something"
						tags = [
							{
								name = "my_tag"
								value_type = "field"
								value = ""
							}
						]
					}`,
				ExpectError: regexp.MustCompile(`Attribute tags\[0\].value string length must be at least 1`),
			},

			// value_field OR value_number is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_event_to_metric_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						metric_name = "my_metric"
						metric_type = "counter"
						metric_kind = "absolute"
						namespace_field = ".namespace"
					}`,
				ExpectError: regexp.MustCompile(`(?s)No attribute specified when one \(and only one\) of.*\[value_field.<.value_number\] is required`),
			},

			// value_field and value_number are mutually exclusive
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_event_to_metric_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						metric_name = "my_metric"
						metric_type = "counter"
						metric_kind = "absolute"
						namespace_field = ".namespace"
						value_field = ".something"
						value_number = 3587
					}`,
				ExpectError: regexp.MustCompile(`(?s)2 attributes specified when one \(and only one\) of.*\[value_field.<.value_number\] is required`),
			},

			// namespace_field and namespace_value are mutually exclusive
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_event_to_metric_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						metric_name = "my_metric"
						metric_type = "counter"
						metric_kind = "absolute"
						namespace_field = ".namespace"
						namespace_value = "nope"
						value_field = ".something"
					}`,
				ExpectError: regexp.MustCompile(`(?s)Attribute "namespace_value" cannot be specified when "namespace_field" is.*specified`),
			},

			// Create
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_event_to_metric_processor" "my_processor" {
						title = "title"
						description = "desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						metric_name = "my_metric"
						metric_type = "counter"
						metric_kind = "absolute"
						namespace_field = ".namespace"
						value_field = ".something"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_event_to_metric_processor.my_processor", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_event_to_metric_processor.my_processor", map[string]any{
						"pipeline_id":     "#mezmo_pipeline.test_parent.id",
						"title":           "title",
						"description":     "desc",
						"generation_id":   "0",
						"inputs.#":        "0",
						"metric_kind":     "absolute",
						"metric_name":     "my_metric",
						"metric_type":     "counter",
						"namespace_field": ".namespace",
						"value_field":     ".something",
						"%":               "14",
					}),
				),
			},

			// Update
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_event_to_metric_processor" "my_processor" {
						title = "updated title"
						description = "updated desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						metric_name = "my_metric"
						metric_type = "gauge"
						metric_kind = "incremental"
						namespace_value = "my_namespace"
						value_number = 55.67
						tags = [
							{
								name = "my_tag"
								value_type = "value"
								value = "tag_value"
							}
						]
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_event_to_metric_processor.my_processor", map[string]any{
						"pipeline_id":     "#mezmo_pipeline.test_parent.id",
						"title":           "updated title",
						"description":     "updated desc",
						"generation_id":   "1",
						"inputs.#":        "1",
						"inputs.0":        "#mezmo_http_source.my_source.id",
						"metric_kind":     "incremental",
						"metric_name":     "my_metric",
						"metric_type":     "gauge",
						"namespace_value": "my_namespace",
						"value_number":    "55.67",
						"%":               "14",
					}),
				),
			},

			// Update - nuke tags and namespace
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_event_to_metric_processor" "my_processor" {
						title = "updated title"
						description = "updated desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						metric_name = "my_metric"
						metric_type = "gauge"
						metric_kind = "incremental"
						value_number = 55
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_event_to_metric_processor.my_processor", map[string]any{
						"%":             "14",
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "updated title",
						"description":   "updated desc",
						"generation_id": "2",
						"inputs.#":      "1",
						"inputs.0":      "#mezmo_http_source.my_source.id",
						"metric_kind":   "incremental",
						"metric_name":   "my_metric",
						"metric_type":   "gauge",
						"value_number":  "55",
					}),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
