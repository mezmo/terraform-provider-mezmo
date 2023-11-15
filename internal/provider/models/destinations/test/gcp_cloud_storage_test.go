package destinations

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestGcpCloudStorageSinkResource(t *testing.T) {
	const cacheKey = "gcp_cloud_storage_destination_resource"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Required fields
			{
				Config: GetProviderConfig() + `
					resource "mezmo_gcp_cloud_storage_destination" "my_dest" {
						inputs = ["abc"]
						pipeline_id = "pipeline-id"
						auth = {
							type = "api_key"
							value = "key"
						}
					}`,
				ExpectError: regexp.MustCompile("The argument \"bucket\" is required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_gcp_cloud_storage_destination" "my_dest" {
						inputs = ["abc"]
						pipeline_id = "pipeline-id"
						bucket = "test_bucket"
					}`,
				ExpectError: regexp.MustCompile("The argument \"auth\" is required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_gcp_cloud_storage_destination" "my_dest" {
						inputs = ["abc"]
						pipeline_id = "pipeline-id"
						bucket = "test_bucket"
						auth = {
							value = "key"
						}
					}`,
				ExpectError: regexp.MustCompile("attribute \"type\" is required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_gcp_cloud_storage_destination" "my_dest" {
						inputs = ["abc"]
						pipeline_id = "pipeline-id"
						bucket = "test_bucket"
						auth = {
							type = "api_key"
						}
					}`,
				ExpectError: regexp.MustCompile("attribute \"value\" is required"),
			},
			// validators
			{
				Config: GetProviderConfig() + `
					resource "mezmo_gcp_cloud_storage_destination" "my_dest" {
						inputs = ["abc"]
						pipeline_id = "pipeline-id"
						bucket = "test_bucket"
						encoding = "invalid"
						auth = {
							type = "api_key"
							value = "key"
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute encoding value must be one of"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_gcp_cloud_storage_destination" "my_dest" {
						inputs = ["abc"]
						pipeline_id = "pipeline-id"
						bucket = ""
						auth = {
							type = "api_key"
							value = "key"
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute bucket string length must be at least 1"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_gcp_cloud_storage_destination" "my_dest" {
						inputs = ["abc"]
						pipeline_id = "pipeline-id"
						bucket = "test_bucket"
						compression = "bzip"
						auth = {
							type = "api_key"
							value = "key"
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute compression value must be one of"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_gcp_cloud_storage_destination" "my_dest" {
						inputs = ["abc"]
						pipeline_id = "pipeline-id"
						bucket = "test_bucket"
						bucket_prefix = ""
						auth = {
							type = "api_key"
							value = "key"
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute bucket_prefix string length must be at least 1"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_gcp_cloud_storage_destination" "my_dest" {
						inputs = ["abc"]
						pipeline_id = "pipeline-id"
						bucket = "test_bucket"
						auth = {
							type = "invalid"
							value = "key"
						}
					}
				`,
				ExpectError: regexp.MustCompile("Attribute auth.type value must be one of"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_gcp_cloud_storage_destination" "my_dest" {
						inputs = ["abc"]
						pipeline_id = "pipeline-id"
						bucket = "test_bucket"
						auth = {
							type = "api_key"
							value = ""
						}
					}
				`,
				ExpectError: regexp.MustCompile("Attribute auth.value string length must be at least 1"),
			},
			// Create
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_http_source" "test_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`) + `
					resource "mezmo_gcp_cloud_storage_destination" "my_dest" {
						title = "test dest"
						description = "test dest description"
						inputs = [mezmo_http_source.test_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id
						encoding = "json"
						compression = "gzip"
						bucket = "test_bucket"
						bucket_prefix = "bucket_prefix"
						auth = {
							type = "api_key"
							value = "key"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_gcp_cloud_storage_destination.my_dest", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_gcp_cloud_storage_destination.my_dest", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "test dest",
						"description":   "test dest description",
						"generation_id": "0",
						"ack_enabled":   "true",
						"inputs.#":      "1",
						"encoding":      "json",
						"compression":   "gzip",
						"bucket":        "test_bucket",
						"bucket_prefix": "bucket_prefix",
						"auth.type":     "api_key",
						"auth.value":    "key",
					}),
				),
			},
			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_gcp_cloud_storage_destination" "import_target" {
						title = "test dest"
						description = "test dest description"
						inputs = [mezmo_http_source.test_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id
						encoding = "json"
						compression = "gzip"
						bucket = "test_bucket"
						bucket_prefix = "bucket_prefix"
						auth = {
							type = "api_key"
							value = "key"
						}
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_gcp_cloud_storage_destination.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_gcp_cloud_storage_destination.my_dest"),
				ImportStateVerify: true,
			},
			// Update fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_http_source" "new_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}
					resource "mezmo_gcp_cloud_storage_destination" "my_dest" {
						title = "new test dest"
						description = "this is a new test description"
						inputs = [mezmo_http_source.new_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id
						encoding = "text"
						compression = "none"
						bucket = "new_bucket"
						bucket_prefix = "newprefix"
						ack_enabled = false
						auth = {
							type = "credentials_json"
							value = "{}"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_gcp_cloud_storage_destination.my_dest", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_gcp_cloud_storage_destination.my_dest", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "new test dest",
						"description":   "this is a new test description",
						"generation_id": "1",
						"ack_enabled":   "false",
						"inputs.#":      "1",
						"encoding":      "text",
						"compression":   "none",
						"bucket":        "new_bucket",
						"bucket_prefix": "newprefix",
						"auth.type":     "credentials_json",
						"auth.value":    "{}",
					}),
				),
			},
		},
	})
}
