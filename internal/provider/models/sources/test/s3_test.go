package sources

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/providertest"
)

func TestS3SourceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Required field "sqs_queue_url"
			{
				Config: GetProviderConfig() + `
					resource "mezmo_s3_source" "my_source" {
						pipeline = "c5ce0dae-0c40-11ee-be56-0242ac120002"
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
						pipeline = "c5ce0dae-0c40-11ee-be56-0242ac120002"
						region = "us-east-1"
						sqs_queue_url = "https://hello.com/sqs"
					}`,
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			// Create and Read testing
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_s3_source" "my_source" {
						pipeline = mezmo_pipeline.test_parent.id
						region = "us-east-2"
						sqs_queue_url = "https://hello.com/sqs"
						auth = {
							access_key_id = "123"
							secret_access_key = "secret123"
						}
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify user-defined properties
					resource.TestCheckResourceAttr("mezmo_s3_source.my_source", "sqs_queue_url", "https://hello.com/sqs"),
					resource.TestCheckResourceAttr("mezmo_s3_source.my_source", "region", "us-east-2"),
					// Verify computed properties
					resource.TestCheckResourceAttrSet("mezmo_s3_source.my_source", "id"),
					resource.TestCheckResourceAttr("mezmo_s3_source.my_source", "generation_id", "0"),
				),
			},
			// Update and Read testing
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_s3_source" "my_source" {
						pipeline = mezmo_pipeline.test_parent.id
						region = "us-east-2"
						sqs_queue_url = "https://hello.com/sqs2"
						auth = {
							access_key_id = "123"
							secret_access_key = "secret123"
						}
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mezmo_s3_source.my_source", "sqs_queue_url", "https://hello.com/sqs2"),
					resource.TestCheckResourceAttr("mezmo_s3_source.my_source", "generation_id", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
