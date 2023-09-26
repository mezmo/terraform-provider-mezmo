package destinations

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestLokiDestinationResource(t *testing.T) {
	cacheKey := "loki_destination_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Check defaults
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title       = "pipeline"
					}
					resource "mezmo_http_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`) + `
					resource "mezmo_loki_destination" "my_destination" {
						title = "test destination"
						description = "loki destination"
						pipeline_id = mezmo_pipeline.test_parent.id
						auth = {
							strategy = "basic"
							user     = "username"
							password = "secret-password"
						}
						endpoint = "http://example.com"
						encoding = "json"
						labels = {
							"test_key_0" = "test_value_0"
							"test_key_1" = "test_value_1"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_loki_destination.my_destination", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_loki_destination.my_destination", map[string]any{
						"title":             "test destination",
						"description":       "loki destination",
						"pipeline_id":       "#mezmo_pipeline.test_parent.id",
						"generation_id":     "0",
						"ack_enabled":       "true",
						"inputs.#":          "0",
						"endpoint":          "http://example.com",
						"encoding":          "json",
						"auth.strategy":     "basic",
						"auth.user":         "username",
						"auth.password":     "secret-password",
						"path":              "/loki/api/v1/push",
						"labels.test_key_0": "test_value_0",
						"labels.test_key_1": "test_value_1",
					}),
				),
			},

			// Required field encoding
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_loki_destination" "my_destination" {
						title = "test destination"
						description = "loki destination"
						pipeline_id = mezmo_pipeline.test_parent.id
						auth = {
							strategy = "basic"
							user     = "username"
							password = "secret-password"
						}
						encoding = "json"
						path     = "this/path"
						labels = {
							"test_key_0" = "test_value_0"
							"test_key_1" = "test_value_1"
						}
					}`,
				ExpectError: regexp.MustCompile("The argument \"endpoint\" is required, but no definition was found"),
			},

			// Required field encoding
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_loki_destination" "my_destination" {
						title = "test destination"
						description = "loki destination"
						pipeline_id = mezmo_pipeline.test_parent.id
						auth = {
							strategy = "basic"
							user     = "username"
							password = "secret-password"
						}
						endpoint = "http://test-endpoint.com"
						path     = "this/path"
						labels = {
							"test_key_0" = "test_value_0"
							"test_key_1" = "test_value_1"
						}
					}`,
				ExpectError: regexp.MustCompile("The argument \"encoding\" is required, but no definition was found"),
			},

			// Required field auth
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_loki_destination" "my_destination" {
						title = "test destination"
						description = "loki destination"
						pipeline_id = mezmo_pipeline.test_parent.id
						endpoint = "http://test-endpoint.com"
						encoding = "json"
						path     = "this/path"
						labels = {
							"test_key_0" = "test_value_0"
							"test_key_1" = "test_value_1"
						}
					}`,
				ExpectError: regexp.MustCompile("The argument \"auth\" is required, but no definition was found"),
			},

			// Required field labels
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_loki_destination" "my_destination" {
						title = "test destination"
						description = "loki destination"
						pipeline_id = mezmo_pipeline.test_parent.id
						endpoint = "http://test-endpoint.com"
						encoding = "json"
						path     = "this/path"
					}`,
				ExpectError: regexp.MustCompile("The argument \"auth\" is required, but no definition was found"),
			},

			// Update all fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_loki_destination" "my_destination" {
						title = "new test destination"
						description = "new loki destination"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						auth = {
							strategy = "basic"
							user     = "new-username"
							password = "new-secret-password"
						}
						endpoint = "http://newexample.com"
						encoding = "text"
						path     = "new/path"
						labels = {
							"new_test_key_0" = "new_test_value_0"
							"new_test_key_1" = "new_test_value_1"
						}
						ack_enabled = false
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_loki_destination.my_destination", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_loki_destination.my_destination", map[string]any{
						"title":                 "new test destination",
						"description":           "new loki destination",
						"pipeline_id":           "#mezmo_pipeline.test_parent.id",
						"generation_id":         "1",
						"ack_enabled":           "false",
						"inputs.#":              "1",
						"endpoint":              "http://newexample.com",
						"encoding":              "text",
						"auth.strategy":         "basic",
						"auth.user":             "new-username",
						"auth.password":         "new-secret-password",
						"path":                  "new/path",
						"labels.new_test_key_0": "new_test_value_0",
						"labels.new_test_key_1": "new_test_value_1",
					}),
				),
			},
		},
	})
}
