package sources

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestWebhookSource(t *testing.T) {
	cacheKey := "webhook_source_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Error: Required field "pipeline_id"
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
							title = "parent pipeline"
						}`) + `
					resource "mezmo_webhook_source" "my_source" {
						signing_key = "nope"
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},
			// Error: "signing_key" length requirements
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_webhook_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						signing_key = ""
					}`,
				ExpectError: regexp.MustCompile("Attribute signing_key string length must be at least 1"),
			},
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_webhook_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						signing_key = "` + strings.Repeat("x", 513) + `"
					}`,
				ExpectError: regexp.MustCompile("Attribute signing_key string length must be at most 512"),
			},

			// Create and Read testing with defaults
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_webhook_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						description = "my description"
						title = "my title"
						signing_key = "sshhh"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("mezmo_webhook_source.my_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestMatchResourceAttr("mezmo_webhook_source.my_source", "gateway_route_id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_webhook_source.my_source", map[string]any{
						"description":      "my description",
						"title":            "my title",
						"generation_id":    "0",
						"signing_key":      "sshhh",
						"capture_metadata": "false",
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_webhook_source" "import_target" {
						pipeline_id = mezmo_pipeline.test_parent.id
						description = "my description"
						title = "my title"
						signing_key = "sshhh"
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_webhook_source.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_webhook_source.my_source"),
				ImportStateVerify: true,
			},

			// Updates
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_webhook_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						description = "new description"
						title = "new title"
						signing_key = "updated"
						capture_metadata = "true"
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_webhook_source.my_source", map[string]any{
						"description":      "new description",
						"title":            "new title",
						"generation_id":    "1",
						"signing_key":      "updated",
						"capture_metadata": "true",
					}),
				),
			},

			// Supply gateway_route_id
			// TODO: We need to take signing_key out of user_config. Currently, you cannot "share" a
			// TODO: source simply by providing a gateway_route_id. We should not expose this sharing to customers as of yet.
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_webhook_source" "parent_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my title"
						description = "my description"
						signing_key = "key1"
					}`) + `
					resource "mezmo_webhook_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "A shared source"
						description = "This source provides gateway_route_id"
						gateway_route_id = mezmo_webhook_source.parent_source.gateway_route_id
						signing_key = mezmo_webhook_source.parent_source.signing_key
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_webhook_source.shared_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_webhook_source.shared_source", map[string]any{
						"description":      "This source provides gateway_route_id",
						"generation_id":    "0",
						"title":            "A shared source",
						"signing_key":      "key1",
						"capture_metadata": "false",
						"pipeline_id":      "#mezmo_pipeline.test_parent.id",
						"gateway_route_id": "#mezmo_webhook_source.parent_source.gateway_route_id",
					}),
				),
			},

			// Updating gateway_route_id is not allowed
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_webhook_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						gateway_route_id = mezmo_pipeline.test_parent.id
						signing_key = mezmo_webhook_source.parent_source.signing_key
					}`,
				ExpectError: regexp.MustCompile("This field is immutable after resource creation."),
			},

			// gateway_route_id can be specified if it's the same value
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_webhook_source" "shared_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "Updated title"
						gateway_route_id = mezmo_webhook_source.parent_source.gateway_route_id
						signing_key = mezmo_webhook_source.parent_source.signing_key
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_webhook_source.shared_source", map[string]any{
						"title":            "Updated title",
						"generation_id":    "1",
						"gateway_route_id": "#mezmo_webhook_source.parent_source.gateway_route_id",
					}),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
