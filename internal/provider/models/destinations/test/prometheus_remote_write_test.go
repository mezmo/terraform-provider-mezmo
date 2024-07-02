package destinations

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestAccPrometheusDestinationResource(t *testing.T) {
	const cacheKey = "prometheus_destination_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Error: pipeline_id is required
			{
				Config: GetProviderConfig() + `
					resource "mezmo_prometheus_remote_write_destination" "my_destination" {
						endpoint = "http://example.com"
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},

			// Error: uri is required
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_prometheus_remote_write_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile("The argument \"endpoint\" is required"),
			},

			// Error: auth.strategy is required
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_prometheus_remote_write_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						endpoint = "https://example.com"
						auth = {}
					}`,
				ExpectError: regexp.MustCompile("attribute \"strategy\" is required"),
			},

			// Error: basic auth fields required
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_prometheus_remote_write_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						endpoint = "https://example.com"
						auth = {
							strategy = "basic"
						}
					}`,
				ExpectError: regexp.MustCompile("Basic auth requires user and password fields to be defined"),
			},

			// Error: basic auth fields required
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_prometheus_remote_write_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						endpoint = "https://example.com"
						auth = {
							strategy = "bearer"
						}
					}`,
				ExpectError: regexp.MustCompile("Bearer auth requires token field to be defined"),
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
					resource "mezmo_prometheus_remote_write_destination" "my_destination" {
						title = "My destination"
						description = "my destination description"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						endpoint = "https://google.com"
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_prometheus_remote_write_destination.my_destination", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_prometheus_remote_write_destination.my_destination", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "My destination",
						"description":   "my destination description",
						"generation_id": "0",
						"ack_enabled":   "true",
						"inputs.#":      "1",
						"inputs.0":      "#mezmo_http_source.my_source.id",
						"endpoint":      "https://google.com",
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_prometheus_remote_write_destination" "import_target" {
						title = "My destination"
						description = "my destination description"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						endpoint = "https://google.com"
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_prometheus_remote_write_destination.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_prometheus_remote_write_destination.my_destination"),
				ImportStateVerify: true,
			},

			// Update all fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_prometheus_remote_write_destination" "my_destination" {
						title = "My destination"
						description = "my destination description"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						endpoint = "https://google.com"
						ack_enabled = false
						auth = {
							strategy = "basic"
							user = "my_user"
							password = "my_password"
						}
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_prometheus_remote_write_destination.my_destination", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_prometheus_remote_write_destination.my_destination", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "My destination",
						"description":   "my destination description",
						"generation_id": "1",
						"ack_enabled":   "false",
						"inputs.#":      "1",
						"inputs.0":      "#mezmo_http_source.my_source.id",
						"endpoint":      "https://google.com",
						"auth.strategy": "basic",
						"auth.user":     "my_user",
						"auth.password": "my_password",
					}),
				),
			},

			// Update auth
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_prometheus_remote_write_destination" "my_destination" {
						title = "My destination"
						description = "my destination description"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						endpoint = "https://google.com"
						auth = {
							strategy = "bearer"
							token = "my_token"
						}
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_prometheus_remote_write_destination.my_destination", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_prometheus_remote_write_destination.my_destination", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "My destination",
						"description":   "my destination description",
						"generation_id": "2",
						"ack_enabled":   "true",
						"inputs.#":      "1",
						"inputs.0":      "#mezmo_http_source.my_source.id",
						"endpoint":      "https://google.com",
						"auth.strategy": "bearer",
						"auth.token":    "my_token",
					}),
				),
			},

			// API-level validation
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_prometheus_remote_write_destination" "my_destination" {
						title = "My destination"
						description = "my destination description"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						endpoint = "not a valid url"
					}
					`,
				ExpectError: regexp.MustCompile("(?s)endpoint.*Must be a valid.*URI"),
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
				resource "mezmo_prometheus_remote_write_destination" "test_destination" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					inputs 			= [mezmo_http_source.my_source2.id]
					endpoint = "https://google.com"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_prometheus_remote_write_destination.test_destination", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_prometheus_remote_write_destination.test_destination", "title", "new title"),
					// verify resource will be re-created after refresh
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_prometheus_remote_write_destination.test_destination",
					),
				),
				ExpectNonEmptyPlan: true,
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
