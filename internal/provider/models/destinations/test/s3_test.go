package destinations

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestAccS3DestinationResource(t *testing.T) {
	const cacheKey = "s3_destination_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Error: properties are required
			{
				Config: GetProviderConfig() + `
					resource "mezmo_s3_destination" "my_destination" {
						region = "us-east2"
					}`,
				ExpectError: regexp.MustCompile("The argument \"auth\" is required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_s3_destination" "my_destination" {
						region = "us-east2"
						auth = {
							access_key_id = "my_key"
						}
					}`,
				ExpectError: regexp.MustCompile("(?s)attribute \"secret_access_key\".*is.*required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_s3_destination" "my_destination" {
						region = "us-east2"
						auth = {
							secret_access_key = "my_secret"
						}
					}`,
				ExpectError: regexp.MustCompile("(?s)attribute \"access_key_id\".*is.*required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_s3_destination" "my_destination" {
						auth = {
							access_key_id = "my_key"
							secret_access_key = "my_secret"
						}
					}`,
				ExpectError: regexp.MustCompile("The argument \"region\" is required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_s3_destination" "my_destination" {
						region = "eu-south-2"
						auth = {
							access_key_id = "my_key"
							secret_access_key = "my_secret"
						}
					}`,
				ExpectError: regexp.MustCompile("The argument \"bucket\" is required"),
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
					resource "mezmo_s3_destination" "my_destination" {
						title       = "My destination"
						description = "my destination description"
						pipeline_id = mezmo_pipeline.test_parent.id
						region      = "us-west1"
						auth = {
							access_key_id = "my_key"
							secret_access_key = "my_secret"
						}
						bucket = "mybucket"
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_s3_destination.my_destination", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_s3_destination.my_destination", map[string]any{
						"pipeline_id":            "#mezmo_pipeline.test_parent.id",
						"title":                  "My destination",
						"description":            "my destination description",
						"generation_id":          "0",
						"ack_enabled":            "true",
						"batch_timeout_secs":     "300",
						"inputs.#":               "0",
						"region":                 "us-west1",
						"auth.access_key_id":     "my_key",
						"auth.secret_access_key": "my_secret",
						"compression":            "none",
						"encoding":               "text",
						"prefix":                 "/",
						"bucket":                 "mybucket",
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_s3_destination" "import_target" {
						title       = "My destination"
						description = "my destination description"
						pipeline_id = mezmo_pipeline.test_parent.id
						region      = "us-west1"
						auth = {
							access_key_id = "my_key"
							secret_access_key = "my_secret"
						}
						bucket = "mybucket"
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_s3_destination.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_s3_destination.my_destination"),
				ImportStateVerify: true,
			},

			// Update all fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_s3_destination" "my_destination" {
						title              = "My new title"
						description        = "My new description"
						inputs             = [mezmo_http_source.my_source.id]
						pipeline_id        = mezmo_pipeline.test_parent.id
						ack_enabled        = false
						batch_timeout_secs = 30
						region             = "us-west2"
						auth = {
							access_key_id = "my_key2"
							secret_access_key = "my_secret2"
						}
						bucket = "mybucket2"
						prefix = "/abc/"
						compression = "gzip"
						encoding = "json"
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_s3_destination.my_destination", map[string]any{
						"pipeline_id":            "#mezmo_pipeline.test_parent.id",
						"title":                  "My new title",
						"description":            "My new description",
						"generation_id":          "1",
						"ack_enabled":            "false",
						"batch_timeout_secs":     "30",
						"inputs.#":               "1",
						"region":                 "us-west2",
						"auth.access_key_id":     "my_key2",
						"auth.secret_access_key": "my_secret2",
						"compression":            "gzip",
						"encoding":               "json",
						"prefix":                 "/abc/",
						"bucket":                 "mybucket2",
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
				resource "mezmo_s3_destination" "test_destination" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					inputs 			= [mezmo_http_source.my_source2.id]
					region      = "us-west1"
					auth = {
						access_key_id = "my_key"
						secret_access_key = "my_secret"
					}
					bucket = "mybucket"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_s3_destination.test_destination", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_s3_destination.test_destination", "title", "new title"),
					// verify resource will be re-created after refresh
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_s3_destination.test_destination",
					),
				),
				ExpectNonEmptyPlan: true,
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
