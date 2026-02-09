package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/v5/internal/provider/providertest"
)

func TestAccessKeyResource(t *testing.T) {
	cacheKey := "access_key_resource_tests"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Cache base resources of pipeline and source
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "my_pipeline" {
						title = "pipeline"
					}
					resource "mezmo_http_source" "my_source" {
						pipeline_id = mezmo_pipeline.my_pipeline.id
					}`,
				),
			},
			// Create testing
			{
				Config: GetCachedConfig(cacheKey) + `
					output "http_key" {
						value = mezmo_access_key.for_http.key
						sensitive = true
						description = "the output of key"
					}
					resource "mezmo_access_key" "for_http" {
						title = "http ingestion key"
						source_id = mezmo_http_source.my_source.id
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchOutput("http_key", regexp.MustCompile(`\w+`)),
					resource.TestMatchResourceAttr(
						"mezmo_access_key.for_http", "id", regexp.MustCompile(`[\w-]{36}`),
					),
					StateHasExpectedValues("mezmo_access_key.for_http", map[string]any{
						"title":     "http ingestion key",
						"source_id": "#mezmo_http_source.my_source.id",
					}),
				),
			},
			// Access keys are immutable, so no updates are allowed
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_access_key" "for_http" {
						title = "change not allowed"
						source_id = mezmo_http_source.my_source.id
					}`,
				ExpectError: regexp.MustCompile(
					`(?s).*Cannot update mezmo_access_key after creation.*` +
						`Access keys are currently imutable and cannot be updated`,
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
