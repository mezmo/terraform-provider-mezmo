package destinations

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestAccGcpCloudPubSubSinkResource(t *testing.T) {
	const cacheKey = "gcp_cloud_pubsub_destination_resource"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Required fields
			{
				Config: GetProviderConfig() + `
					resource "mezmo_gcp_cloud_pubsub_destination" "my_dest" {
						inputs = ["abc"]
						pipeline_id = "pipeline-id"
						credentials_json = "{}"
					}`,
				ExpectError: regexp.MustCompile("The argument \"topic\" is required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_gcp_cloud_pubsub_destination" "my_dest" {
						inputs = ["abc"]
						pipeline_id = "pipeline-id"
						topic = "topic12"
					}`,
				ExpectError: regexp.MustCompile("The argument \"credentials_json\" is required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_gcp_cloud_pubsub_destination" "my_dest" {
						inputs = ["abc"]
						pipeline_id = "pipeline-id"
						topic = "topic12"
						credentials_json = "{}"
					}`,
				ExpectError: regexp.MustCompile("The argument \"project_id\" is required"),
			},
			// validators
			{
				Config: GetProviderConfig() + `
					resource "mezmo_gcp_cloud_pubsub_destination" "my_dest" {
						inputs = ["abc"]
						pipeline_id = "pipeline-id"
						bucket = "test_bucket"
						encoding = "invalid"
						credentials_json = "{}"
					}`,
				ExpectError: regexp.MustCompile("Attribute encoding value must be one of"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_gcp_cloud_pubsub_destination" "my_dest" {
						inputs = ["abc"]
						pipeline_id = "pipeline-id"
						topic = ""
						credentials_json = "{}"
					}`,
				ExpectError: regexp.MustCompile("Attribute topic string length must be at least 1"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_gcp_cloud_pubsub_destination" "my_dest" {
						inputs = ["abc"]
						pipeline_id = "pipeline-id"
						bucket = "test_bucket"
						topic = "topic"
						project_id = ""
						credentials_json = "{}"
					}`,
				ExpectError: regexp.MustCompile("Attribute project_id string length must be at least 1"),
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
					resource "mezmo_gcp_cloud_pubsub_destination" "my_dest" {
						title = "test dest"
						description = "test dest description"
						inputs = [mezmo_http_source.test_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id
						encoding = "json"
						project_id = "proj1"
						topic = "topic1"
						credentials_json = "{}"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_gcp_cloud_pubsub_destination.my_dest", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_gcp_cloud_pubsub_destination.my_dest", map[string]any{
						"pipeline_id":      "#mezmo_pipeline.test_parent.id",
						"title":            "test dest",
						"description":      "test dest description",
						"generation_id":    "0",
						"ack_enabled":      "true",
						"inputs.#":         "1",
						"encoding":         "json",
						"project_id":       "proj1",
						"topic":            "topic1",
						"credentials_json": "{}",
					}),
				),
			},
			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_gcp_cloud_pubsub_destination" "import_target" {
						title = "test dest"
						description = "test dest description"
						inputs = [mezmo_http_source.test_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id
						encoding = "json"
						project_id = "proj1"
						topic = "topic1"
						credentials_json = "{}"
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_gcp_cloud_pubsub_destination.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_gcp_cloud_pubsub_destination.my_dest"),
				ImportStateVerify: true,
			},
			// Update fields, remove bucket prefix since it's optional
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_http_source" "new_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}
					resource "mezmo_gcp_cloud_pubsub_destination" "my_dest" {
						title = "new test dest"
						description = "this is a new test description"
						inputs = [mezmo_http_source.new_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id
						encoding = "text"
						project_id = "proj2"
						topic = "topic2"
						ack_enabled = false
						credentials_json = "{}"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_gcp_cloud_pubsub_destination.my_dest", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_gcp_cloud_pubsub_destination.my_dest", map[string]any{
						"pipeline_id":      "#mezmo_pipeline.test_parent.id",
						"title":            "new test dest",
						"description":      "this is a new test description",
						"generation_id":    "1",
						"ack_enabled":      "false",
						"inputs.#":         "1",
						"encoding":         "text",
						"topic":            "topic2",
						"project_id":       "proj2",
						"credentials_json": "{}",
					}),
				),
			},
			// confirm manually deleted resources are recreated
			{
				Config: GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_http_source" "my_source2" {
					pipeline_id = mezmo_pipeline.test_parent2.id
				}
				resource "mezmo_gcp_cloud_pubsub_destination" "test_destination" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					encoding 		= "json"
					project_id = "proj1"
					topic = "topic1"
					credentials_json = "{}"
					inputs 			= [mezmo_http_source.my_source2.id]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_gcp_cloud_pubsub_destination.test_destination", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_gcp_cloud_pubsub_destination.test_destination", "title", "new title"),
					// verify resource will be re-created after refresh
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_gcp_cloud_pubsub_destination.test_destination",
					),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}