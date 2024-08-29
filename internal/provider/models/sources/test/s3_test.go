package sources

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/providertest"
)

func TestAccS3SourceResource(t *testing.T) {
	const cacheKey = "s3_source_resource"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Required field "sqs_queue_url"
			{
				Config: GetProviderConfig() + `
					resource "mezmo_s3_source" "my_source" {
						pipeline_id = "c5ce0dae-0c40-11ee-be56-0242ac120002"
						region = "us-east-2"
						auth = {
							access_key_id = "123"
							secret_access_key = "secret123"
						}
					}`,
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			// Required field "auth"
			{
				Config: GetProviderConfig() + `
					resource "mezmo_s3_source" "my_source" {
						pipeline_id = "c5ce0dae-0c40-11ee-be56-0242ac120002"
						region = "us-east-1"
						sqs_queue_url = "https://hello.com/sqs"
					}`,
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			// Violation of compression enum
			{
				Config: GetProviderConfig() + `
					resource "mezmo_s3_source" "my_source" {
						pipeline_id = "c5ce0dae-0c40-11ee-be56-0242ac120002"
						region = "us-east-1"
						sqs_queue_url = "https://hello.com/sqs"
						auth = {
							access_key_id = "123"
							secret_access_key = "secret123"
						}
						compression = "NOPE"
					}`,
				ExpectError: regexp.MustCompile("Attribute compression value must be one of:"),
			},
			// Create and Read testing
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}`) + `
					resource "mezmo_s3_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my title"
						description = "my description"
						region = "us-east-2"
						sqs_queue_url = "https://hello.com/sqs"
						auth = {
							access_key_id = "123"
							secret_access_key = "secret123"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("mezmo_s3_source.my_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_s3_source.my_source", map[string]any{
						"description":            "my description",
						"title":                  "my title",
						"region":                 "us-east-2",
						"sqs_queue_url":          "https://hello.com/sqs",
						"auth.access_key_id":     "123",
						"auth.secret_access_key": "secret123",
						"generation_id":          "0",
						"compression":            "auto",
					}),
				),
			},
			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_s3_source" "import_target" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my title"
						description = "my description"
						region = "us-east-2"
						sqs_queue_url = "https://hello.com/sqs"
						auth = {
							access_key_id = "123"
							secret_access_key = "secret123"
						}
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_s3_source.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_s3_source.my_source"),
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_s3_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "new title"
						description = "new desc"
						region = "us-east-1"
						sqs_queue_url = "https://hello.com/sqs2"
						auth = {
							access_key_id = "4321"
							secret_access_key = "4321secret"
						}
						compression = "gzip"
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_s3_source.my_source", map[string]any{
						"description":            "new desc",
						"title":                  "new title",
						"region":                 "us-east-1",
						"sqs_queue_url":          "https://hello.com/sqs2",
						"auth.access_key_id":     "4321",
						"auth.secret_access_key": "4321secret",
						"generation_id":          "1",
						"compression":            "gzip",
					}),
				),
			},
			// confirm manually deleted resources are recreated
			{
				Config: GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_s3_source" "test_source" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					region 			= "us-east-2"
					sqs_queue_url = "https://hello.com/sqs"
					auth 					= {
						access_key_id = "123"
						secret_access_key = "secret123"
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_s3_source.test_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_s3_source.test_source", "title", "new title"),
					// delete the resource
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_s3_source.test_source",
					),
				),
				// verify resource will be re-created after refresh
				ExpectNonEmptyPlan: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
