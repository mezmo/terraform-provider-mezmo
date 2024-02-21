package destinations

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestMezmoDestinationResource(t *testing.T) {
	const cacheKey = "mezmo_destination_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Error: pipeline_id is required
			{
				Config: GetProviderConfig() + `
					resource "mezmo_logs_destination" "my_destination" {
						uri = "http://example.com"
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},
			// Error: incorrect options
			{
				Config: GetProviderConfig() + `
					resource "mezmo_logs_destination" "my_destination" {
						pipeline_id = "pip1"
						ingestion_key = "my_key"
						log_construction_scheme = "pass-through"
						explicit_scheme_options = {
							line = "abc"
						}
					}`,
				ExpectError: regexp.MustCompile("Invalid log constructions options"),
			},

			// Create test defaults
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_http_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`) + `
					resource "mezmo_logs_destination" "my_destination" {
						title = "My destination"
						description = "my destination description"
						pipeline_id = mezmo_pipeline.test_parent.id
						ingestion_key = "my_key"
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_logs_destination.my_destination", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_logs_destination.my_destination", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "My destination",
						"description":   "my destination description",
						"generation_id": "0",
						"ack_enabled":   "true",
						"inputs.#":      "0",
						"host":          "logs.logdna.com",
						"ingestion_key": "my_key",
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_logs_destination" "import_target" {
						title = "My destination"
						description = "my destination description"
						pipeline_id = mezmo_pipeline.test_parent.id
						ingestion_key = "my_key"
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_logs_destination.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_logs_destination.my_destination"),
				ImportStateVerify: true,
			},

			// Update all fields (pass through)
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_logs_destination" "my_destination" {
						title = "new title"
						description = "new description"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						ack_enabled = "false"
						host = "zzz.mezmo.com"
						ingestion_key = "key2"
						log_construction_scheme = "pass-through"
						query = {
							hostname = "{{ .host }}",
							ip = "{{metadata.query.ip}}",
							mac = "{{metadata.query.mac}}"
							tags = ["tag1"]
						}
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_logs_destination.my_destination", map[string]any{
						"pipeline_id":             "#mezmo_pipeline.test_parent.id",
						"title":                   "new title",
						"description":             "new description",
						"generation_id":           "1",
						"inputs.#":                "1",
						"inputs.0":                "#mezmo_http_source.my_source.id",
						"ack_enabled":             "false",
						"host":                    "zzz.mezmo.com",
						"ingestion_key":           "key2",
						"log_construction_scheme": "pass-through",
						"query.hostname":          "{{ .host }}",
						"query.mac":               "{{metadata.query.mac}}",
						"query.ip":                "{{metadata.query.ip}}",
						"query.tags.#":            "1",
						"query.tags.0":            "tag1",
					}),
				),
			},

			// Update explicit fields and nullify others
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_logs_destination" "my_destination" {
						title = "new title"
						description = "new description"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						ingestion_key = "key3"
						log_construction_scheme = "explicit"
						explicit_scheme_options = {
							line       = ".thing_one",
							app        = "{{metadata.query.app}}",
							meta_field = "metadata._meta",
							env        = "{{._env}}"
						}
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_logs_destination.my_destination", map[string]any{
						"pipeline_id":                        "#mezmo_pipeline.test_parent.id",
						"title":                              "new title",
						"description":                        "new description",
						"generation_id":                      "2",
						"inputs.#":                           "1",
						"inputs.0":                           "#mezmo_http_source.my_source.id",
						"ack_enabled":                        "true",
						"host":                               "logs.logdna.com",
						"ingestion_key":                      "key3",
						"log_construction_scheme":            "explicit",
						"explicit_scheme_options.line":       ".thing_one",
						"explicit_scheme_options.app":        "{{metadata.query.app}}",
						"explicit_scheme_options.meta_field": "metadata._meta",
						"explicit_scheme_options.env":        "{{._env}}",
					}),
				),
			},

			// API-level validation
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_logs_destination" "my_destination" {
						title         = "My destination"
						description   = "my destination description"
						pipeline_id   = mezmo_pipeline.test_parent.id
						inputs        = [mezmo_http_source.my_source.id]
						ingestion_key = "my_key"
						host          = "invalid.host.com"
					}
					`,
				ExpectError: regexp.MustCompile("match pattern"),
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
				resource "mezmo_logs_destination" "test_destination" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					inputs 			= [mezmo_http_source.my_source2.id]
					ingestion_key = "my_key"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_logs_destination.test_destination", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_logs_destination.test_destination", "title", "new title"),
					// verify resource will be re-created after refresh
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_logs_destination.test_destination",
					),
				),
				ExpectNonEmptyPlan: true,
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
