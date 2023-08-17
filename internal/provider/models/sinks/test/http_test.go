package sinks

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/providertest"
)

func TestHttpSinkResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Error: pipeline_id is required
			{
				Config: GetProviderConfig() + `
					resource "mezmo_http_sink" "my_sink" {
						uri = "http://example.com"
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},

			// Error: uri is required
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_http_sink" "my_sink" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile("The argument \"uri\" is required"),
			},

			// Create test defaults
			{
				Config: SetCachedConfig("http_sink_resources", `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_http_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`) + `
					resource "mezmo_http_sink" "my_sink" {
						title = "my http sink"
						description = "http sink description"
						pipeline_id = mezmo_pipeline.test_parent.id
						uri = "http://example.com"
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_http_sink.my_sink", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_http_sink.my_sink", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "my http sink",
						"description":   "http sink description",
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
				Config: GetCachedConfig("http_sink_resources") + `
					resource "mezmo_http_sink" "my_sink" {
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
					StateHasExpectedValues("mezmo_http_sink.my_sink", map[string]any{
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
				Config: GetCachedConfig("http_sink_resources") + `
					resource "mezmo_http_sink" "my_sink" {
						pipeline_id = mezmo_pipeline.test_parent.id
						uri = "http://example.com"
						inputs = [mezmo_http_source.my_source.id]
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_http_sink.my_sink", map[string]any{
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

			// API-level validation
			{
				Config: GetCachedConfig("http_sink_resources") + `
					resource "mezmo_http_sink" "my_sink" {
						pipeline_id = mezmo_pipeline.test_parent.id
						uri = "http://example.com"
						inputs = [mezmo_http_source.my_source.id]

						auth = {
							strategy = "bearer"
							user = "nope"
							password = "nope"
						}
					}
					`,
				ExpectError: regexp.MustCompile("properties/token/minLength"),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
