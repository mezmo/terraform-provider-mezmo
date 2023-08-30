package transforms

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/providertest"
)

func TestReduceTransform(t *testing.T) {
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
					resource "mezmo_reduce_transform" "my_transform" {
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},

			// // Error: `field` is required
			// {
			// 	Config: GetCachedConfig(cacheKey) + `
			// 		resource "mezmo_reduce_transform" "my_transform" {
			// 			pipeline_id = mezmo_pipeline.test_parent.id
			// 		}`,
			// 	ExpectError: regexp.MustCompile("The argument \"field\" is required"),
			// },

			// Create
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_reduce_transform" "my_transform" {
						title = "transform title"
						description = "transform desc"
						pipeline_id = mezmo_pipeline.test_parent.id

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
								]
								logical_operation = "OR"
							}
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_reduce_transform.my_transform", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_reduce_transform.my_transform", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "transform title",
						"description":   "transform desc",
						"generation_id": "0",
						"inputs.#":      "0",
					}),
				),
			},

			// // Update
			// {
			// 	Config: GetCachedConfig(cacheKey) + `
			// 		resource "mezmo_reduce_transform" "my_transform" {
			// 			title = "new title"
			// 			description = "new desc"
			// 			pipeline_id = mezmo_pipeline.test_parent.id
			// 			inputs = [mezmo_http_source.my_source.id]
			// 			field = ".thing2"
			// 			values_only = false
			// 		}`,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		StateHasExpectedValues("mezmo_reduce_transform.my_transform", map[string]any{
			// 			"pipeline_id":   "#mezmo_pipeline.test_parent.id",
			// 			"title":         "new title",
			// 			"description":   "new desc",
			// 			"generation_id": "1",
			// 			"inputs.#":      "1",
			// 			"inputs.0":      "#mezmo_http_source.my_source.id",
			// 			"field":         ".thing2",
			// 			"values_only":   "false",
			// 		}),
			// 	),
			// },

			// // Error: server-side validation
			// {
			// 	Config: GetCachedConfig(cacheKey) + `
			// 	resource "mezmo_reduce_transform" "my_transform" {
			// 		pipeline_id = mezmo_pipeline.test_parent.id
			// 		inputs = []
			// 		field = "not-a-valid-field"
			// 	}`,
			// 	ExpectError: regexp.MustCompile("match pattern"),
			// },

			// Delete testing automatically occurs in TestCase
		},
	})
}
