package transforms

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/providertest"
)

func TestScriptExecutionTransform(t *testing.T) {
	const cacheKey = "script_execution_resources"
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
					resource "mezmo_script_execution_transform" "my_transform" {
						script = "function processEvent(e) { return e }"
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},

			// Error: `field` is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_script_execution_transform" "my_transform" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile("The argument \"script\" is required"),
			},

			// Error: `field` values validates length
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_script_execution_transform" "my_transform" {
						pipeline_id = mezmo_pipeline.test_parent.id
						script = ""
					}`,
				ExpectError: regexp.MustCompile("Invalid Attribute Value Length"),
			},

			// Create
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_script_execution_transform" "my_transform" {
						title = "transform title"
						description = "transform desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						script = "function processEvent(e) { return e }"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_script_execution_transform.my_transform", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_script_execution_transform.my_transform", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "transform title",
						"description":   "transform desc",
						"generation_id": "0",
						"inputs.#":      "0",
						"script":        "function processEvent(e) { return e }",
					}),
				),
			},

			// Update
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_script_execution_transform" "my_transform" {
					title = "new title"
					description = "new desc"
					pipeline_id = mezmo_pipeline.test_parent.id
					inputs = [mezmo_http_source.my_source.id]
					script = <<-EOT
					  function processEvent(e) {
					    if (e.skip) {
					      return null
					    }
					    return e
					  }
					  EOT
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_script_execution_transform.my_transform", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "new title",
						"description":   "new desc",
						"generation_id": "1",
						"inputs.#":      "1",
						"inputs.0":      "#mezmo_http_source.my_source.id",
						"script": `function processEvent(e) {
  if (e.skip) {
    return null
  }
  return e
}
`,
					}),
				),
			},

			// Error: server-side validation
			{
				Config: GetCachedConfig(cacheKey) + `
				resource "mezmo_script_execution_transform" "my_transform" {
					pipeline_id = mezmo_pipeline.test_parent.id
					inputs = [mezmo_http_source.my_source.id]
					script = "THIS IS NOT VALID"
				}`,
				ExpectError: regexp.MustCompile("script is not valid JavaScript"),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
