package destinations

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestHttpDestination(t *testing.T) {
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
						uri = "http://example.com"
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
					}),
				),
			},

			// Update all fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_http_destination" "my_destination" {
						title = "new title"
						description = "new description"
						pipeline_id = mezmo_pipeline.test_parent.id
						uri = "http://example.com"
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
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_http_destination.my_destination", map[string]any{
						"pipeline_id":           "#mezmo_pipeline.test_parent.id",
						"title":                 "new title",
						"description":           "new description",
						"generation_id":         "1",
						"encoding":              "json",
						"compression":           "gzip",
						"ack_enabled":           "false",
						"inputs.#":              "1",
						"inputs.0":              "#mezmo_http_source.my_source.id",
						"auth.strategy":         "basic",
						"auth.user":             "me",
						"auth.password":         "ssshhh",
						"headers.x-my-header":   "first header",
						"headers.x-your-header": "second header",
					}),
				),
			},

			// Nullify fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_http_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						uri = "http://example.com"
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

			// Delete testing automatically occurs in TestCase
		},
	})
}
