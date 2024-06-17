package destinations

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestAccAzureBlobStorageDestinationResource(t *testing.T) {
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

			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_azure_blob_storage_destination" "my_destination" {
						title             = "My destination"
						description       = "my destination description"
						inputs            = [mezmo_http_source.my_source.id]
						pipeline_id       = mezmo_pipeline.test_parent.id
						connection_string = "abc://defg.com"
						container_name    = "my_container"
					}
					resource "mezmo_azure_blob_storage_destination" "implicit" {
						title              = "My destination"
						description        = "my destination description"
						inputs             = [mezmo_http_source.my_source.id]
						pipeline_id        = mezmo_pipeline.test_parent.id
						connection_string  = "foo"
						container_name     = "bar"
						file_consolidation = {
							enabled = true
						}
					}
					resource "mezmo_azure_blob_storage_destination" "explicit" {
						title              = "My destination"
						description        = "my destination description"
						inputs             = [mezmo_http_source.my_source.id]
						pipeline_id        = mezmo_pipeline.test_parent.id
						connection_string  = "foo"
						container_name     = "bar"
						file_consolidation = {
							enabled = true
							process_every_seconds = 420
							requested_size_bytes = 420000000
							base_path = "/foo/bar"
						}
					}
					`,
				Check: resource.ComposeTestCheckFunc([]resource.TestCheckFunc{
					// Not enabled
					resource.TestCheckNoResourceAttr("mezmo_azure_blob_storage_destination.my_destination", "file_consolidation"),
					// Implicit config
					resource.TestCheckResourceAttr("mezmo_azure_blob_storage_destination.implicit", "file_consolidation.enabled", "true"),
					resource.TestCheckResourceAttr("mezmo_azure_blob_storage_destination.implicit", "file_consolidation.process_every_seconds", "600"),
					resource.TestCheckResourceAttr("mezmo_azure_blob_storage_destination.implicit", "file_consolidation.requested_size_bytes", "500000000"),
					resource.TestCheckResourceAttr("mezmo_azure_blob_storage_destination.implicit", "file_consolidation.base_path", ""),
					testAccDestinationBackend("mezmo_azure_blob_storage_destination.implicit", map[string]any{
						"ack_enabled":        true,
						"batch_timeout_secs": float64(300),
						"encoding":           "text",
						"compression":        "none",
						"container_name":     "bar",
						"connection_string":  "foo",
						"file_consolidation": map[string]any{
							"enabled":               true,
							"process_every_seconds": float64(600),
							"requested_size_bytes":  float64(5e+08), // The backend returns this in scientific notation 🤷‍♂️
							"base_path":             "",
						},
					}),
					// Explicit config
					resource.TestCheckResourceAttr("mezmo_azure_blob_storage_destination.explicit", "file_consolidation.enabled", "true"),
					resource.TestCheckResourceAttr("mezmo_azure_blob_storage_destination.explicit", "file_consolidation.process_every_seconds", "420"),
					resource.TestCheckResourceAttr("mezmo_azure_blob_storage_destination.explicit", "file_consolidation.requested_size_bytes", "420000000"),
					resource.TestCheckResourceAttr("mezmo_azure_blob_storage_destination.explicit", "file_consolidation.base_path", "/foo/bar"),
					testAccDestinationBackend("mezmo_azure_blob_storage_destination.explicit", map[string]any{
						"ack_enabled":        true,
						"batch_timeout_secs": float64(300),
						"encoding":           "text",
						"compression":        "none",
						"container_name":     "bar",
						"connection_string":  "foo",
						"file_consolidation": map[string]any{
							"enabled":               true,
							"process_every_seconds": float64(420),
							"requested_size_bytes":  float64(4.2e+08), // The backend returns this in scientific notation 🤷‍♂️
							"base_path":             "/foo/bar",
						},
					}),
				}...),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_azure_blob_storage_destination" "import_target" {
						title             = "My destination"
						description       = "my destination description"
						inputs            = [mezmo_http_source.my_source.id]
						pipeline_id       = mezmo_pipeline.test_parent.id
						connection_string = "abc://defg.com"
						container_name    = "my_container"
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_azure_blob_storage_destination.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_azure_blob_storage_destination.my_destination"),
				ImportStateVerify: true,
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
			// confirm manually deleted resources are recreated
			{
				Config: GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_azure_blob_storage_destination" "new_destination" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					inputs 			= []
					connection_string = "abc://defg.com"
					container_name    = "my_container"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_azure_blob_storage_destination.new_destination", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_azure_blob_storage_destination.new_destination", "title", "new title"),
					// delete the resource
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_azure_blob_storage_destination.new_destination",
					),
				),
				// verify resource will be re-created after refresh
				ExpectNonEmptyPlan: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
