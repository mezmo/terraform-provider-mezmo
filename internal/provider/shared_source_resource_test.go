package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestSharedSourceResource(t *testing.T) {
	cacheKey := "shared_source_tests"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: SetCachedConfig(cacheKey, `
					output "shared_source_key" {
						value = mezmo_access_key.shared.key
						sensitive = true
					}
					resource "mezmo_access_key" "shared" {
						title = "An access key for the shared source"
						source_id = mezmo_shared_source.my_source.id
					}`) + `
					resource "mezmo_shared_source" "my_source" {
						title = "HTTP Shared"
            description = "This source can be shared across pipelines"
            type = "http"
          }
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchOutput("shared_source_key", regexp.MustCompile(`\w+`)),
					resource.TestMatchResourceAttr("mezmo_shared_source.my_source", "id", IDRegex),
					StateHasExpectedValues("mezmo_shared_source.my_source", map[string]any{
						"title":       "HTTP Shared",
						"description": "This source can be shared across pipelines",
						"type":        "http",
					}),
				),
			},
			// Update testing
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_shared_source" "my_source" {
						title = "updated title"
            description = "updated description"
						type = "http"
          }
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchOutput("shared_source_key", regexp.MustCompile(`\w+`)),
					resource.TestMatchResourceAttr("mezmo_shared_source.my_source", "id", IDRegex),
					StateHasExpectedValues("mezmo_shared_source.my_source", map[string]any{
						"title":       "updated title",
						"description": "updated description",
					}),
				),
			},
			// Updating `type` causes the whole resource to be re-created
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_shared_source" "my_source" {
						title = "updated title"
            description = "updated description"
						type = "kinesis-firehose"
          }
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchOutput("shared_source_key", regexp.MustCompile(`\w+`)),
					resource.TestMatchResourceAttr("mezmo_shared_source.my_source", "id", IDRegex),
					resource.TestCheckResourceAttrPair("mezmo_access_key.shared", "source_id", "mezmo_shared_source.my_source", "id"),
					StateHasExpectedValues("mezmo_shared_source.my_source", map[string]any{
						"title":       "updated title",
						"description": "updated description",
						"type":        "kinesis-firehose",
					}),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestSharedSourceResourceErrors(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		ErrorCheck: CheckMultipleErrors([]string{
			`Attribute title string length must be at least 1`,
			`Attribute description string length must be at least 1`,
			`Attribute type value must be one of: \["http" "splunk-hec".*`,
		}),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `
					resource "mezmo_shared_source" "my_source" {
						title = ""
            description = ""
            type = ""
          }
				`,
			},
		},
	})
}
