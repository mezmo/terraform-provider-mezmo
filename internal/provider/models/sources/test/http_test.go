package sources

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/providertest"
)

func TestHttpSourceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Required fields parent pipeline id
			{
				Config: GetProviderConfig() + `
					resource "mezmo_http_source" "my_source" {}`,
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			// Validator tests
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_http_source" "my_source" {
						pipeline = mezmo_pipeline.test_parent.id
						decoding = "nope"
					}`,
				ExpectError: regexp.MustCompile("Attribute decoding value must be one of:"),
			},
			// Create and Read testing
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_http_source" "my_source" {
						pipeline = mezmo_pipeline.test_parent.id
						title = "my http title"
						description = "my http description"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_http_source.my_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_http_source.my_source", map[string]any{
						"description":      "my http description",
						"generation_id":    "0",
						"title":            "my http title",
						"decoding":         "json",
						"capture_metadata": "false",
						"pipeline":         "#mezmo_pipeline.test_parent.id",
					}),
				),
			},
			// Update and Read testing
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_http_source" "my_source" {
						pipeline = mezmo_pipeline.test_parent.id
						title = "new title"
						description = "new description"
						decoding = "ndjson"
						capture_metadata = "true"
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_http_source.my_source", map[string]any{
						"description":      "new description",
						"generation_id":    "1",
						"title":            "new title",
						"decoding":         "ndjson",
						"capture_metadata": "true",
					}),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
