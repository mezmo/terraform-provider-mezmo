package sources

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestAzureEventHubSourceResource(t *testing.T) {
	cacheKey := "azure_event_hub_source_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { providertest.TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Required field test cases
			{
				Config: providertest.GetProviderConfig() + `
					resource "mezmo_azure_event_hub_source" "req_field_source" {
						connection_string = "test_connection_string"
						namespace = "test_namespace"
						group_id = "test_group_id"
						topics = ["test_topic"]
					}`,
				ExpectError: regexp.MustCompile("argument \"pipeline_id\" is required"),
			},
			{
				Config: providertest.GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_azure_event_hub_source" "req_field_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						namespace = "test_namespace"
						group_id = "test_group_id"
						topics = ["test_topic"]
					}`,
				ExpectError: regexp.MustCompile("argument \"connection_string\" is required"),
			},
			{
				Config: providertest.GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_azure_event_hub_source" "req_field_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						connection_string = "test_connection_string"
						group_id = "test_group_id"
						topics = ["test_topic"]
					}`,
				ExpectError: regexp.MustCompile("argument \"namespace\" is required"),
			},
			{
				Config: providertest.GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_azure_event_hub_source" "req_field_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						connection_string = "test_connection_string"
						namespace = "test_namespace"
						topics = ["test_topic"]
					}`,
				ExpectError: regexp.MustCompile("argument \"group_id\" is required"),
			},
			{
				Config: providertest.GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_azure_event_hub_source" "req_field_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						connection_string = "test_connection_string"
						namespace = "test_namespace"
						group_id = "test_group_id"
					}`,
				ExpectError: regexp.MustCompile("argument \"topics\" is required"),
			},
			// Validator test cases
			{
				Config: providertest.GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_azure_event_hub_source" "val_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						decoding = "invalid"
						connection_string = "test_connection_string"
						namespace = "test_namespace"
						group_id = "test_group_id"
						topics = ["test_topic"]
					}`,
				ExpectError: regexp.MustCompile("Attribute decoding value must be one of:"),
			},
			{
				Config: providertest.GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_azure_event_hub_source" "val_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						connection_string = ""
						namespace = "test_namespace"
						group_id = "test_group_id"
						topics = ["test_topic"]
					}`,
				ExpectError: regexp.MustCompile("Attribute connection_string string length must be at least 1"),
			},
			{
				Config: providertest.GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_azure_event_hub_source" "val_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						connection_string = "test_connection_string"
						namespace = ""
						group_id = "test_group_id"
						topics = ["test_topic"]
					}`,
				ExpectError: regexp.MustCompile("Attribute namespace string length must be at least 1"),
			},
			{
				Config: providertest.GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_azure_event_hub_source" "val_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						connection_string = "test_connection_string"
						namespace = "test_namespace"
						group_id = ""
						topics = ["test_topic"]
					}`,
				ExpectError: regexp.MustCompile("Attribute group_id string length must be at least 1"),
			},
			{
				Config: providertest.GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_azure_event_hub_source" "val_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						connection_string = "test_connection_string"
						namespace = "test_namespace"
						group_id = "test_group_id"
						topics = []
					}`,
				ExpectError: regexp.MustCompile("topics list must contain at least 1 elements"),
			},
			{
				Config: providertest.GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_azure_event_hub_source" "val_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						connection_string = "test_connection_string"
						namespace = "test_namespace"
						group_id = "test_group_id"
						topics = [""]
					}`,
				ExpectError: regexp.MustCompile("Attribute topics\\[0] string length must be at least 1"),
			},
			// Create and Validate State
			{
				Config: providertest.SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}`) + `
					resource "mezmo_azure_event_hub_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "test title"
						description = "test description"
						connection_string = "test_connection_string"
						namespace = "test_namespace"
						group_id = "test_group_id"
						topics = ["topic1", "topic2"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_azure_event_hub_source.my_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					providertest.StateHasExpectedValues("mezmo_azure_event_hub_source.my_source", map[string]any{
						"title":             "test title",
						"description":       "test description",
						"generation_id":     "0",
						"decoding":          "bytes",
						"connection_string": "test_connection_string",
						"namespace":         "test_namespace",
						"group_id":          "test_group_id",
						"topics.#":          "2",
						"topics.0":          "topic1",
						"topics.1":          "topic2",
					}),
				),
			},
			{
				Config: providertest.GetCachedConfig(cacheKey) + `
					resource "mezmo_azure_event_hub_source" "import_target" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "test title"
						description = "test description"
						connection_string = "test_connection_string"
						namespace = "test_namespace"
						group_id = "test_group_id"
						topics = ["topic1", "topic2"]
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_azure_event_hub_source.import_target",
				ImportStateIdFunc: providertest.ComputeImportId("mezmo_azure_event_hub_source.my_source"),
				ImportStateVerify: true,
			},
			{
				Config: providertest.GetCachedConfig(cacheKey) + `
					resource "mezmo_azure_event_hub_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "new title"
						description = "new description"
						decoding = "json"
						connection_string = "new_connection_string"
						namespace = "new_namespace"
						group_id = "new_group_id"
						topics = ["topic_new"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					providertest.StateHasExpectedValues("mezmo_azure_event_hub_source.my_source", map[string]any{
						"title":             "new title",
						"description":       "new description",
						"generation_id":     "1",
						"decoding":          "json",
						"connection_string": "new_connection_string",
						"namespace":         "new_namespace",
						"group_id":          "new_group_id",
						"topics.#":          "1",
						"topics.0":          "topic_new",
					}),
				),
			},
			// confirm manually deleted resources are recreated
			{
				Config: providertest.GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_azure_event_hub_source" "test_source" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					connection_string = "test_connection_string"
					namespace 				= "test_namespace"
					group_id 					= "test_group_id"
					topics 						= ["topic1", "topic2"]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_azure_event_hub_source.test_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_azure_event_hub_source.test_source", "title", "new title"),
					// delete the resource
					providertest.TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_azure_event_hub_source.test_source",
					),
				),
				// verify resource will be re-created after refresh
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
