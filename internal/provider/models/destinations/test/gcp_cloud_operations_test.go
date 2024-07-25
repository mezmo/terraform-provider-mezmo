package destinations

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestGcpCloudOperationsSinkResource_errors(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		ErrorCheck: CheckMultipleErrors([]string{
			"The argument \"resource_type\" is required",
			"The argument \"log_id\" is required",
			"The argument \"project_id\" is required",
		}),
		Steps: []resource.TestStep{
			// Required fields
			{
				Config: GetProviderConfig() + `
									resource "mezmo_gcp_cloud_operations_destination" "rf_resource_type" {
										inputs = ["abc"]
										pipeline_id = "pipeline-id"
										credentials_json = "{}"
									}`,
			},
		},
	})
}

func TestGcpCloudOperationsSinkResource_validation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		ErrorCheck: CheckMultipleErrors([]string{
			"Attribute resource_type string length must be at least 1",
			"Attribute project_id string length must be at least 1",
			"Attribute log_id string length must be at least 1",
		}),
		Steps: []resource.TestStep{
			// Required fields
			{
				Config: GetProviderConfig() + `
									resource "mezmo_gcp_cloud_operations_destination" "val3" {
										inputs = ["abc"]
										pipeline_id = "pipeline-id"
										credentials_json = "{}"
										resource_type = ""
										log_id = ""
										project_id = ""
									}`,
			},
		},
	})
}

func TestAccGcpCloudOperationsSinkResource_crud(t *testing.T) {
	const cacheKey = "gcp_cloud_operations_destination_resource"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Create
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_http_source" "test_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`) + `
					resource "mezmo_gcp_cloud_operations_destination" "my_dest" {
						title = "test dest"
						description = "test dest description"
						inputs = [mezmo_http_source.test_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id
						resource_type = "global"
						project_id = "proj123"
						log_id = "log123"
						credentials_json = "{}"
						resource_labels = {
							"somekey" = "someval"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_gcp_cloud_operations_destination.my_dest", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_gcp_cloud_operations_destination.my_dest", map[string]any{
						"pipeline_id":             "#mezmo_pipeline.test_parent.id",
						"title":                   "test dest",
						"description":             "test dest description",
						"generation_id":           "0",
						"ack_enabled":             "true",
						"inputs.#":                "1",
						"resource_type":           "global",
						"project_id":              "proj123",
						"log_id":                  "log123",
						"credentials_json":        "{}",
						"resource_labels.%":       "1",
						"resource_labels.somekey": "someval",
					}),
				),
			},
			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_gcp_cloud_operations_destination" "import_target" {
						title = "test dest"
						description = "test dest description"
						inputs = [mezmo_http_source.test_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id
						resource_type = "global"
						project_id = "proj123"
						log_id = "log123"
						credentials_json = "{}"
						resource_labels = {
							somekey = "someval"
						}
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_gcp_cloud_operations_destination.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_gcp_cloud_operations_destination.my_dest"),
				ImportStateVerify: true,
			},
			// update
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_http_source" "new_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}
					resource "mezmo_gcp_cloud_operations_destination" "my_dest" {
						title = "new test dest"
						description = "this is a new test description"
						inputs = [mezmo_http_source.new_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id
						resource_type = "global2"
						log_id = "log456"
						project_id = "proj456"
						credentials_json = "{}"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_gcp_cloud_operations_destination.my_dest", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_gcp_cloud_operations_destination.my_dest", map[string]any{
						"pipeline_id":       "#mezmo_pipeline.test_parent.id",
						"title":             "new test dest",
						"description":       "this is a new test description",
						"generation_id":     "1",
						"ack_enabled":       "true",
						"inputs.#":          "1",
						"resource_type":     "global2",
						"resource_labels.%": nil,
						"log_id":            "log456",
						"project_id":        "proj456",
						"credentials_json":  "{}",
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
				resource "mezmo_gcp_cloud_operations_destination" "test_destination" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					resource_type = "global2"
					project_id = "proj456"
				  log_id = "log123"
					credentials_json = "{}"
					inputs 			= [mezmo_http_source.my_source2.id]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_gcp_cloud_operations_destination.test_destination", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_gcp_cloud_operations_destination.test_destination", "title", "new title"),
					// verify resource will be re-created after refresh
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_gcp_cloud_operations_destination.test_destination",
					),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
