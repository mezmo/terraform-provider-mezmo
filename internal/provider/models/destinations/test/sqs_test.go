package destinations

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/providertest"
)

func TestAccSQSSinkResource(t *testing.T) {
	const cacheKey = "sqs_sink_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Error: properties are required
			{
				Config: GetProviderConfig() + `
					resource "mezmo_sqs_destination" "my_destination" {
						region = "us-east-1"
					}`,
				ExpectError: regexp.MustCompile("The argument \"auth\" is required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_sqs_destination" "my_destination" {
						region = "us-east-1"
						auth = {
							access_key_id = "my_key"
						}
					}`,
				ExpectError: regexp.MustCompile("(?s)attribute \"secret_access_key\".*is.*required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_sqs_destination" "my_destination" {
						region = "us-east-1"
						auth = {
							secret_access_key = "my_secret"
						}
					}`,
				ExpectError: regexp.MustCompile("(?s)attribute \"access_key_id\".*is.*required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_sqs_destination" "my_destination" {
						auth = {
							access_key_id = "my_key"
							secret_access_key = "my_secret"
						}
					}`,
				ExpectError: regexp.MustCompile("The argument \"region\" is required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_sqs_destination" "my_destination" {
						region = "us-east-1"
						auth = {
							access_key_id = "my_key"
							secret_access_key = "my_secret"
						}
					}`,
				ExpectError: regexp.MustCompile("The argument \"queue_url\" is required"),
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
					resource "mezmo_sqs_destination" "my_destination" {
						title       = "My destination"
						description = "my destination description"
						pipeline_id = mezmo_pipeline.test_parent.id
						region      = "us-east-1"
						auth = {
							access_key_id = "my_key"
							secret_access_key = "my_secret"
						}
						queue_url = "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_sqs_destination.my_destination", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_sqs_destination.my_destination", map[string]any{
						"pipeline_id":            "#mezmo_pipeline.test_parent.id",
						"title":                  "My destination",
						"description":            "my destination description",
						"generation_id":          "0",
						"ack_enabled":            "true",
						"batch_timeout_secs":     "300",
						"inputs.#":               "0",
						"region":                 "us-east-1",
						"auth.access_key_id":     "my_key",
						"auth.secret_access_key": "my_secret",
						"compression":            "none",
						"encoding":               "text",
						"queue_url":              "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
					}),
				),
			},

			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_sqs_destination" "my_destination" {
						title       = "My destination"
						description = "my destination description"
						pipeline_id = mezmo_pipeline.test_parent.id
						region      = "us-east-1"
						auth = {
							access_key_id = "my_key"
							secret_access_key = "my_secret"
						}
						queue_url = "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"
					}
					resource "mezmo_sqs_destination" "implicit" {
						title             = "My destination"
						description       = "my destination description"
						inputs            = [mezmo_http_source.my_source.id]
						pipeline_id       = mezmo_pipeline.test_parent.id
						region            = "us-east-1"
						auth              = {
							access_key_id     = "my_key"
							secret_access_key = "my_secret"
						}
						queue_url = "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"
					}
					resource "mezmo_sqs_destination" "explicit" {
						title             = "My destination"
						description       = "my destination description"
						inputs            = [mezmo_http_source.my_source.id]
						pipeline_id       = mezmo_pipeline.test_parent.id
						region            = "us-east-1"
						auth              = {
							access_key_id     = "my_key"
							secret_access_key = "my_secret"
						}
						queue_url = "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"
					}
					`,
				Check: resource.ComposeTestCheckFunc([]resource.TestCheckFunc{
					// Not enabled
					resource.TestCheckNoResourceAttr("mezmo_sqs_destination.my_destination", "file_consolidation"),
					// Implicit config
					resource.TestCheckResourceAttr("mezmo_sqs_destination.implicit", "file_consolidation.enabled", "true"),
					resource.TestCheckResourceAttr("mezmo_sqs_destination.implicit", "file_consolidation.process_every_seconds", "600"),
					resource.TestCheckResourceAttr("mezmo_sqs_destination.implicit", "file_consolidation.requested_size_bytes", "500000000"),
					resource.TestCheckResourceAttr("mezmo_sqs_destination.implicit", "file_consolidation.base_path", ""),
					// Explicit config
					resource.TestCheckResourceAttr("mezmo_sqs_destination.explicit", "file_consolidation.enabled", "true"),
					resource.TestCheckResourceAttr("mezmo_sqs_destination.explicit", "file_consolidation.process_every_seconds", "420"),
					resource.TestCheckResourceAttr("mezmo_sqs_destination.explicit", "file_consolidation.requested_size_bytes", "420000000"),
					resource.TestCheckResourceAttr("mezmo_sqs_destination.explicit", "file_consolidation.base_path", "/foo/bar"),
				}...),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_sqs_destination" "import_target" {
						title       = "My destination"
						description = "my destination description"
						pipeline_id = mezmo_pipeline.test_parent.id
						region      = "us-east-1"
						auth = {
							access_key_id = "my_key"
							secret_access_key = "my_secret"
						}
						queue_url = "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_sqs_destination.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_sqs_destination.my_destination"),
				ImportStateVerify: true,
			},

			// Update all fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_sqs_destination" "my_destination" {
						title              = "My new title"
						description        = "My new description"
						inputs             = [mezmo_http_source.my_source.id]
						pipeline_id        = mezmo_pipeline.test_parent.id
						ack_enabled        = false
						batch_timeout_secs = 30
						region             = "us-east-2"
						auth = {
							access_key_id = "my_key2"
							secret_access_key = "my_secret2"
						}
						queue_url = "https://sqs.us-east-2.amazonaws.com/123456789012/test-queue"
						compression = "gzip"
						encoding = "json"
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_sqs_destination.my_destination", map[string]any{
						"pipeline_id":            "#mezmo_pipeline.test_parent.id",
						"title":                  "My new title",
						"description":            "My new description",
						"generation_id":          "1",
						"ack_enabled":            "false",
						"batch_timeout_secs":     "30",
						"inputs.#":               "1",
						"region":                 "us-east-2",
						"auth.access_key_id":     "my_key2",
						"auth.secret_access_key": "my_secret2",
						"compression":            "gzip",
						"encoding":               "json",
						"queue_url":              "https://sqs.us-east-2.amazonaws.com/123456789012/test-queue",
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
				resource "mezmo_sqs_destination" "test_destination" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					inputs 			= [mezmo_http_source.my_source2.id]
					region      = "us-east-1"
					auth = {
						access_key_id = "my_key"
						secret_access_key = "my_secret"
					}
					queue_url = "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_sqs_destination.test_destination", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_sqs_destination.test_destination", "title", "new title"),
					// verify resource will be re-created after refresh
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_sqs_destination.test_destination",
					),
				),
				ExpectNonEmptyPlan: true,
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
