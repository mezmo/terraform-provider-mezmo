package sources

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/v6/internal/provider/providertest"
)

func TestAccHttpClientSource(t *testing.T) {
	cacheKey := "http_client_source_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Required field pipeline_id
			{
				Config: GetProviderConfig() + `
					resource "mezmo_http_client_source" "my_source" {
						endpoint = "https://example.com/data"
						method = "GET"
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},
			// Required field endpoint
			{
				Config: GetProviderConfig() + `
					resource "mezmo_http_client_source" "my_source" {
						pipeline_id = "c5ce0dae-0c40-11ee-be56-0242ac120002"
						method = "GET"
					}`,
				ExpectError: regexp.MustCompile("The argument \"endpoint\" is required"),
			},
			// Required field method
			{
				Config: GetProviderConfig() + `
					resource "mezmo_http_client_source" "my_source" {
						pipeline_id = "c5ce0dae-0c40-11ee-be56-0242ac120002"
						endpoint = "https://example.com/data"
					}`,
				ExpectError: regexp.MustCompile("The argument \"method\" is required"),
			},
			// Validator: invalid method
			{
				Config: GetProviderConfig() + `
					resource "mezmo_http_client_source" "my_source" {
						pipeline_id = "c5ce0dae-0c40-11ee-be56-0242ac120002"
						endpoint = "https://example.com/data"
						method = "NOPE"
					}`,
				ExpectError: regexp.MustCompile("Attribute method value must be one of:"),
			},
			// Validator: invalid endpoint URL
			{
				Config: GetProviderConfig() + `
					resource "mezmo_http_client_source" "my_source" {
						pipeline_id = "c5ce0dae-0c40-11ee-be56-0242ac120002"
						endpoint = "not a url"
						method = "GET"
					}`,
				ExpectError: regexp.MustCompile("must be a valid http or https URL"),
			},
			// Validator: invalid decoding
			{
				Config: GetProviderConfig() + `
					resource "mezmo_http_client_source" "my_source" {
						pipeline_id = "c5ce0dae-0c40-11ee-be56-0242ac120002"
						endpoint = "https://example.com/data"
						method = "GET"
						decoding = "nope"
					}`,
				ExpectError: regexp.MustCompile("Attribute decoding value must be one of:"),
			},
			// Validator: scrape_interval_secs too low
			{
				Config: GetProviderConfig() + `
					resource "mezmo_http_client_source" "my_source" {
						pipeline_id = "c5ce0dae-0c40-11ee-be56-0242ac120002"
						endpoint = "https://example.com/data"
						method = "GET"
						scrape_interval_secs = 5
					}`,
				ExpectError: regexp.MustCompile("Attribute scrape_interval_secs value must be between 15 and 86400"),
			},
			// Validator: scrape_timeout_secs too low
			{
				Config: GetProviderConfig() + `
					resource "mezmo_http_client_source" "my_source" {
						pipeline_id = "c5ce0dae-0c40-11ee-be56-0242ac120002"
						endpoint = "https://example.com/data"
						method = "GET"
						scrape_timeout_secs = 2
					}`,
				ExpectError: regexp.MustCompile("Attribute scrape_timeout_secs value must be between 5 and 60"),
			},
			// Validator: invalid auth strategy
			{
				Config: GetProviderConfig() + `
					resource "mezmo_http_client_source" "my_source" {
						pipeline_id = "c5ce0dae-0c40-11ee-be56-0242ac120002"
						endpoint = "https://example.com/data"
						method = "GET"
						auth = {
							strategy = "oauth"
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute auth.strategy value must be one of:"),
			},
			// Create and Read testing with defaults
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}`) + `
					resource "mezmo_http_client_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my http client title"
						description = "my http client description"
						endpoint = "https://example.com/data"
						method = "GET"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_http_client_source.my_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_http_client_source.my_source", map[string]any{
						"description":          "my http client description",
						"generation_id":        "0",
						"title":                "my http client title",
						"endpoint":             "https://example.com/data",
						"method":               "GET",
						"scrape_interval_secs": "30",
						"scrape_timeout_secs":  "10",
						"decoding":             "bytes",
						"pipeline_id":          "#mezmo_pipeline.test_parent.id",
					}),
				),
			},
			// Error: basic auth requires user and password
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_http_client_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						endpoint = "https://example.com/data"
						method = "GET"
						auth = {
							strategy = "basic"
						}
					}`,
				ExpectError: regexp.MustCompile("Basic auth requires user and password fields to be defined"),
			},
			// Error: bearer auth requires token
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_http_client_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						endpoint = "https://example.com/data"
						method = "GET"
						auth = {
							strategy = "bearer"
						}
					}`,
				ExpectError: regexp.MustCompile("Bearer auth requires token field to be defined"),
			},
			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_http_client_source" "import_target" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my http client title"
						description = "my http client description"
						endpoint = "https://example.com/data"
						method = "GET"
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_http_client_source.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_http_client_source.my_source"),
				ImportStateVerify: true,
			},
			// Update all fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_http_client_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "new title"
						description = "new description"
						endpoint = "https://example.com/v2/data"
						method = "POST"
						scrape_interval_secs = 60
						scrape_timeout_secs = 30
						decoding = "json"
						auth = {
							strategy = "basic"
							user = "myuser"
							password = "mypassword"
						}
						headers = {
							"X-Custom" = "custom-value"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_http_client_source.my_source", map[string]any{
						"description":          "new description",
						"generation_id":        "1",
						"title":                "new title",
						"endpoint":             "https://example.com/v2/data",
						"method":               "POST",
						"scrape_interval_secs": "60",
						"scrape_timeout_secs":  "30",
						"decoding":             "json",
						"auth.strategy":        "basic",
						"auth.user":            "myuser",
						"auth.password":        "mypassword",
						"headers.X-Custom":     "custom-value",
					}),
				),
			},
			// Update auth to bearer
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_http_client_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "new title"
						description = "new description"
						endpoint = "https://example.com/v2/data"
						method = "GET"
						auth = {
							strategy = "bearer"
							token = "my-secret-token"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_http_client_source.my_source", map[string]any{
						"generation_id": "2",
						"auth.strategy": "bearer",
						"auth.token":    "my-secret-token",
						"method":        "GET",
					}),
				),
			},
			// confirm manually deleted resources are recreated
			{
				Config: GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_http_client_source" "test_source" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title = "new title"
					endpoint = "https://example.com/data"
					method = "GET"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_http_client_source.test_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_http_client_source.test_source", "title", "new title"),
					// delete the resource
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_http_client_source.test_source",
					),
				),
				// verify resource will be re-created after refresh
				ExpectNonEmptyPlan: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
