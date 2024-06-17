package destinations

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestAccHttpDestination(t *testing.T) {
	cacheKey := "http_destination_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Error: pipeline_id is required
			{
				Config: GetProviderConfig() + `
					resource "mezmo_http_destination" "my_destination" {
						uri = "http://example.com"
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},

			// Error: uri is required
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}`) + `
					resource "mezmo_http_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile("The argument \"uri\" is required"),
			},

			// Error: basic auth fields required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_http_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						uri = "http://example.com"
						auth = {
							strategy = "basic"
						}
					}`,
				ExpectError: regexp.MustCompile("Basic auth requires user and password fields to be defined"),
			},

			// Error: basic auth fields required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_http_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						uri = "https://example.com"
						auth = {
							strategy = "bearer"
						}
					}`,
				ExpectError: regexp.MustCompile("Bearer auth requires token field to be defined"),
			},

			// Error: must set both payload_prefix and payload_suffix
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_http_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						uri = "https://example.com"
						payload_prefix = "{\"extra_prop\": true"
					}`,
				ExpectError: regexp.MustCompile("If 'payload_prefix' is set, 'payload_suffix' must be as well."),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_http_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						uri = "https://example.com"
						payload_suffix = "\"extra_prop2\": true }"
					}`,
				ExpectError: regexp.MustCompile("If 'payload_suffix' is set, 'payload_prefix' must be as well."),
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
					resource "mezmo_http_destination" "my_destination" {
						title = "my http destination"
						description = "http destination description"
						pipeline_id = mezmo_pipeline.test_parent.id
						uri = "https://google.com"
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_http_destination.my_destination", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_http_destination.my_destination", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "my http destination",
						"description":   "http destination description",
						"generation_id": "0",
						"encoding":      "text",
						"compression":   "none",
						"ack_enabled":   "true",
						"inputs.#":      "0",
						"uri":           "https://google.com",
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_http_destination" "import_target" {
						title = "my http destination"
						description = "http destination description"
						pipeline_id = mezmo_pipeline.test_parent.id
						uri = "https://example.com"
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_http_destination.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_http_destination.my_destination"),
				ImportStateVerify: true,
			},

			// Update all fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_http_destination" "my_destination" {
						title = "new title"
						description = "new description"
						pipeline_id = mezmo_pipeline.test_parent.id
						uri = "https://google.com"
						encoding = "json"
						compression = "gzip"
						ack_enabled = "false"
						inputs = [mezmo_http_source.my_source.id]
						auth = {
							strategy = "basic"
							user = "me"
							password = "ssshhh"
						}

						headers = {
							"x-my-header" = "first header"
							"x-your-header" = "second header"
						}
						max_bytes = 5000
						timeout_secs = 600
						method = "post"
						payload_prefix = "{\"extra_prop\": true"
						payload_suffix = "\"extra_prop\": true }"
						tls_protocols = ["h2"]
						proxy = {
							enabled = true
							endpoint = "http://myproxy.com"
							hosts_bypass_proxy = ["0.0.0.0", "1.1.1.1"]
						}
						rate_limiting = {
							request_limit = 600
							duration_secs = 900
						}
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_http_destination.my_destination", map[string]any{
						"pipeline_id":                 "#mezmo_pipeline.test_parent.id",
						"title":                       "new title",
						"description":                 "new description",
						"generation_id":               "1",
						"encoding":                    "json",
						"compression":                 "gzip",
						"ack_enabled":                 "false",
						"inputs.#":                    "1",
						"inputs.0":                    "#mezmo_http_source.my_source.id",
						"auth.strategy":               "basic",
						"auth.user":                   "me",
						"auth.password":               "ssshhh",
						"headers.x-my-header":         "first header",
						"headers.x-your-header":       "second header",
						"max_bytes":                   "5000",
						"timeout_secs":                "600",
						"method":                      "post",
						"payload_prefix":              "{\"extra_prop\": true",
						"payload_suffix":              "\"extra_prop\": true }",
						"tls_protocols.#":             "1",
						"tls_protocols.0":             "h2",
						"proxy.enabled":               "true",
						"proxy.endpoint":              "http://myproxy.com",
						"proxy.hosts_bypass_proxy.#":  "2",
						"proxy.hosts_bypass_proxy.0":  "0.0.0.0",
						"proxy.hosts_bypass_proxy.1":  "1.1.1.1",
						"rate_limiting.request_limit": "600",
						"rate_limiting.duration_secs": "900",
					}),
				),
			},

			// Nullify fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_http_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						uri = "https://google.com"
						inputs = [mezmo_http_source.my_source.id]
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_http_destination.my_destination", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"generation_id": "2",
						"encoding":      "text",
						"compression":   "none",
						"ack_enabled":   "true",
						"auth.%":        nil,
						"headers.%":     nil,
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
				resource "mezmo_http_destination" "test_destination" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					uri = "https://google.com"
					inputs 			= [mezmo_http_source.my_source2.id]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_http_destination.test_destination", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_http_destination.test_destination", "title", "new title"),
					// verify resource will be re-created after refresh
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_http_destination.test_destination",
					),
				),
				ExpectNonEmptyPlan: true,
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
