package processors

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/providertest"
)

func TestAccDataProfilerProcessor(t *testing.T) {
	var sb strings.Builder
	for i := 0; i < 512; i++ {
		sb.WriteString("a")
	}
	const cacheKey = "data_profiler_resourcess"
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
					resource "mezmo_data_profiler_processor" "my_processor" {
						app_fields = [".app", ".container"]
  						host_fields = [".host", ".hostname"]
  						level_fields = [".level", ".log_level"]
  						line_fields = [".line", ".message"]
  						label_fields = [".labels"]
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},

			// Error: `app_fields` is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_data_profiler_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
  						host_fields = [".host", ".hostname"]
  						level_fields = [".level", ".log_level"]
  						line_fields = [".line", ".message"]
  						label_fields = [".labels"]
					}`,
				ExpectError: regexp.MustCompile("The argument \"app_fields\" is required"),
			},

			// Error: `host_fields` is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_data_profiler_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						app_fields = [".app", ".container"]
  						level_fields = [".level", ".log_level"]
  						line_fields = [".line", ".message"]
  						label_fields = [".labels"]
					}`,
				ExpectError: regexp.MustCompile("The argument \"host_fields\" is required"),
			},

			// Error: `level_fields` is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_data_profiler_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						app_fields = [".app", ".container"]
  						host_fields = [".host", ".hostname"]
  						line_fields = [".line", ".message"]
  						label_fields = [".labels"]
					}`,
				ExpectError: regexp.MustCompile("The argument \"level_fields\" is required"),
			},

			// Error: `line_fields` is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_data_profiler_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						app_fields = [".app", ".container"]
  						host_fields = [".host", ".hostname"]
  						level_fields = [".level", ".log_level"]
  						label_fields = [".labels"]
					}`,
				ExpectError: regexp.MustCompile("The argument \"line_fields\" is required"),
			},

			// Error: `app_fields` validates min number of items
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_data_profiler_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						app_fields = []
  						host_fields = [".host", ".hostname"]
  						level_fields = [".level", ".log_level"]
  						line_fields = [".line", ".message"]
  						label_fields = [".labels"]
					}`,
				ExpectError: regexp.MustCompile("Attribute app_fields list must contain at least 1 elements, got: 0"),
			},

			// Error: `app_fields` validates max number of items
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_data_profiler_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						app_fields = [".a", ".b", ".c", ".d", ".e", ".f", ".g", ".h", ".i", ".j", ".k"]
  						host_fields = [".host", ".hostname"]
  						level_fields = [".level", ".log_level"]
  						line_fields = [".line", ".message"]
  						label_fields = [".labels"]
					}`,
				ExpectError: regexp.MustCompile("Attribute app_fields list must contain at most 10 elements, got: 11"),
			},

			// Error: `app_fields` validates min item length
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_data_profiler_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						app_fields = [""]
  						host_fields = [".host", ".hostname"]
  						level_fields = [".level", ".log_level"]
  						line_fields = [".line", ".message"]
  						label_fields = [".labels"]
					}`,
				ExpectError: regexp.MustCompile("Attribute app_fields\\[0\\] string length must be at least 1, got: 0"),
			},

			// Error: `app_fields` validates max item length
			{
				Config: GetCachedConfig(cacheKey) + fmt.Sprintf(`
					resource "mezmo_data_profiler_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						app_fields = [".%s"]
  						host_fields = [".host", ".hostname"]
  						level_fields = [".level", ".log_level"]
  						line_fields = [".line", ".message"]
  						label_fields = [".labels"]
					}`, sb.String()),
				ExpectError: regexp.MustCompile("Attribute app_fields\\[0\\] string length must be at most 512, got: 513"),
			},

			// Error: `app_fields` validates max item length
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_data_profiler_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						app_fields = ["not-valid-field-reference"]
  						host_fields = [".host", ".hostname"]
  						level_fields = [".level", ".log_level"]
  						line_fields = [".line", ".message"]
  						label_fields = [".labels"]
					}`,
				ExpectError: regexp.MustCompile("Must be a valid data access syntax"),
			},

			// Success
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_data_profiler_processor" "my_processor" {
						title = "profile!"
						description = "all the things!"
						pipeline_id = mezmo_pipeline.test_parent.id
						app_fields = [".app"]
  						host_fields = [".host", ".hostname"]
  						level_fields = [".level", ".log_level"]
  						line_fields = [".line", ".message"]
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_data_profiler_processor.my_processor", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_data_profiler_processor.my_processor", map[string]any{
						"pipeline_id":    "#mezmo_pipeline.test_parent.id",
						"title":          "profile!",
						"description":    "all the things!",
						"generation_id":  "0",
						"inputs.#":       "0",
						"app_fields.#":   "1",
						"app_fields.0":   ".app",
						"host_fields.#":  "2",
						"host_fields.0":  ".host",
						"host_fields.1":  ".hostname",
						"level_fields.#": "2",
						"level_fields.0": ".level",
						"level_fields.1": ".log_level",
						"line_fields.#":  "2",
						"line_fields.0":  ".line",
						"line_fields.1":  ".message",
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_data_profiler_processor" "import_target" {
						title = "profile!"
						description = "all the things!"
						pipeline_id = mezmo_pipeline.test_parent.id
						app_fields = [".app"]
  						host_fields = [".host", ".hostname"]
  						level_fields = [".level", ".log_level"]
  						line_fields = [".line", ".message"]
  						label_fields = [".labels"]
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_data_profiler_processor.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_data_profiler_processor.my_processor"),
				ImportStateVerify: true,
			},

			// Update fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_data_profiler_processor" "my_processor" {
						title = "profile! again!"
						description = "all the things! even more!"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						app_fields = [".appname"]
  						host_fields = [".hostname"]
  						level_fields = [".nested.level"]
  						line_fields = [".nested.message"]
  						label_fields = [".label_locations"]
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_data_profiler_processor.my_processor", map[string]any{
						"pipeline_id":    "#mezmo_pipeline.test_parent.id",
						"title":          "profile! again!",
						"description":    "all the things! even more!",
						"generation_id":  "1",
						"inputs.#":       "1",
						"app_fields.#":   "1",
						"app_fields.0":   ".appname",
						"host_fields.#":  "1",
						"host_fields.0":  ".hostname",
						"level_fields.#": "1",
						"level_fields.0": ".nested.level",
						"line_fields.#":  "1",
						"line_fields.0":  ".nested.message",
						"label_fields.#": "1",
						"label_fields.0": ".label_locations",
					}),
				),
			},

			// confirm manually deleted resources are recreated
			{
				Config: GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_http_source" "my_source2" {
					pipeline_id = mezmo_pipeline.test_parent2.id
				}
				resource "mezmo_data_profiler_processor" "test_processor" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					inputs = [mezmo_http_source.my_source2.id]
					app_fields = [".appname"]
  					host_fields = [".hostname"]
  					level_fields = [".nested.level"]
  					line_fields = [".nested.message"]
  					label_fields = [".label_locations"]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_data_profiler_processor.test_processor", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_data_profiler_processor.test_processor", "title", "new title"),
					// delete the resource
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_data_profiler_processor.test_processor",
					),
				),
				// verify resource will be re-created after refresh
				ExpectNonEmptyPlan: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
