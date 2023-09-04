package destinations

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/providertest"
)

func TestAzureBlobStorageDestinationResource(t *testing.T) {
	const cacheKey = "azure_blob_destination_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Error: properties are required
			{
				Config: GetProviderConfig() + `
					resource "mezmo_azure_blob_storage_destination" "my_destination" {
						inputs      = ["abc"]
						pipeline_id = "pip1"
						container_name = "my_container"
					}`,
				ExpectError: regexp.MustCompile("The argument \"connection_string\" is required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_azure_blob_storage_destination" "my_destination" {
						inputs      = ["abc"]
						pipeline_id = "pip1"
						connection_string = "my_connection_string"
					}`,
				ExpectError: regexp.MustCompile("The argument \"container_name\" is required"),
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
					resource "mezmo_azure_blob_storage_destination" "my_destination" {
						title             = "My destination"
						description       = "my destination description"
						inputs            = [mezmo_http_source.my_source.id]
						pipeline_id       = mezmo_pipeline.test_parent.id
						connection_string = "abc://defg.com"
						container_name    = "my_container"
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_azure_blob_storage_destination.my_destination", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_azure_blob_storage_destination.my_destination", map[string]any{
						"pipeline_id":        "#mezmo_pipeline.test_parent.id",
						"title":              "My destination",
						"description":        "my destination description",
						"generation_id":      "0",
						"ack_enabled":        "true",
						"inputs.#":           "1",
						"connection_string":  "abc://defg.com",
						"container_name":     "my_container",
						"compression":        "none",
						"encoding":           "text",
						"batch_timeout_secs": "300",
					}),
				),
			},

			// Update all fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_azure_blob_storage_destination" "my_destination" {
						title              = "My destination"
						description        = "my destination description"
						inputs             = [mezmo_http_source.my_source.id]
						pipeline_id        = mezmo_pipeline.test_parent.id
						connection_string  = "abc://zzz.com"
						container_name     = "my_container2"
						compression        = "gzip"
						encoding           = "json"
						prefix             = "a1"
						batch_timeout_secs = 60
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_azure_blob_storage_destination.my_destination", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_azure_blob_storage_destination.my_destination", map[string]any{
						"pipeline_id":        "#mezmo_pipeline.test_parent.id",
						"title":              "My destination",
						"description":        "my destination description",
						"generation_id":      "1",
						"ack_enabled":        "true",
						"inputs.#":           "1",
						"connection_string":  "abc://zzz.com",
						"container_name":     "my_container2",
						"compression":        "gzip",
						"encoding":           "json",
						"prefix":             "a1",
						"batch_timeout_secs": "60",
					}),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
