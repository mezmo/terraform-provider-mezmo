package sinks

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/providertest"
)

func TestElasticSearchSinkResource(t *testing.T) {
	const cacheKey = "elasticsearch_sink_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Error: properties are required
			{
				Config: GetProviderConfig() + `
					resource "mezmo_elasticsearch_sink" "my_sink" {
						inputs      = ["abc"]
						endpoints = ["http://example.com"]
					}`,
				ExpectError: regexp.MustCompile("The argument \"auth\" is required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_elasticsearch_sink" "my_sink" {
						inputs      = ["abc"]
						auth = {
							strategy = "basic"
						}
					}`,
				ExpectError: regexp.MustCompile("The argument \"endpoints\" is required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_elasticsearch_sink" "my_sink" {
						inputs      = ["abc"]
						pipeline_id = "pip1"
						endpoints = ["http://example.com"]
						auth = {
							strategy = "basic"
						}
					}`,
				ExpectError: regexp.MustCompile("Basic auth requires user and password fields to be defined"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_elasticsearch_sink" "my_sink" {
						inputs      = ["abc"]
						pipeline_id = "pip1"
						endpoints = ["http://example.com"]
						auth = {
							strategy = "aws"
						}
					}`,
				ExpectError: regexp.MustCompile("(?s)AWS auth requires access_key_id, secret_access_key and region.*to.*be.*defined"),
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
					resource "mezmo_elasticsearch_sink" "my_sink" {
						title = "my sink"
						description = "my sink description"
						inputs      = [mezmo_http_source.my_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id
						endpoints   = ["https://example.com"]
						auth = {
							strategy = "basic"
							user     = "user1"
							password = "pass1"
						}
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_elasticsearch_sink.my_sink", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_elasticsearch_sink.my_sink", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "my sink",
						"description":   "my sink description",
						"generation_id": "0",
						"ack_enabled":   "true",
						"inputs.#":      "1",
						"endpoints.#":   "1",
						"endpoints.0":   "https://example.com",
						"auth.strategy": "basic",
						"auth.user":     "user1",
						"auth.password": "pass1",
						"compression":   "none",
					}),
				),
			},

			// Update all fields (pass through)
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_http_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`) + `
					resource "mezmo_elasticsearch_sink" "my_sink" {
						title = "new sink"
						description = "new sink description"
						inputs      = [mezmo_http_source.my_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id
						endpoints   = ["https://example2.com"]
						auth = {
							strategy          = "aws"
							region            = "us-east2"
							access_key_id     = "acc1"
							secret_access_key = "secret1"
						}
						index = "another-%Y.%m.%d"
						pipeline = "pip1"
						compression = "gzip"
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_elasticsearch_sink.my_sink", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_elasticsearch_sink.my_sink", map[string]any{
						"pipeline_id":            "#mezmo_pipeline.test_parent.id",
						"title":                  "new sink",
						"description":            "new sink description",
						"generation_id":          "1",
						"ack_enabled":            "true",
						"inputs.#":               "1",
						"endpoints.#":            "1",
						"endpoints.0":            "https://example2.com",
						"auth.strategy":          "aws",
						"auth.region":            "us-east2",
						"auth.access_key_id":     "acc1",
						"auth.secret_access_key": "secret1",
						"index":                  "another-%Y.%m.%d",
						"pipeline":               "pip1",
						"compression":            "gzip",
					}),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
