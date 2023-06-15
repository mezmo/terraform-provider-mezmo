package sources

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/providertest"
)

func TestS3SourceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// // Required fields json
			// {
			// 	Config: GetProviderConfig() + `
			// 		resource "mezmo_pipeline" "test_parent" {
			// 			title = "parent pipeline"
			// 		}
			// 		resource "mezmo_demo_source" "my_source" {
			// 			pipeline = mezmo_pipeline.test_parent.id
			// 		}`,
			// 	ExpectError: regexp.MustCompile("Missing required argument"),
			// },
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
				// // Verify user-defined properties
				// resource.TestCheckResourceAttr("mezmo_demo_source.my_source", "format", "json"),
				// // Verify computed properties
				// resource.TestCheckResourceAttrSet("mezmo_demo_source.my_source", "id"),
				// resource.TestCheckResourceAttrSet("mezmo_demo_source.my_source", "generation_id"),
				),
			},
			// Update and Read testing
			// {
			// 	Config: GetProviderConfig() + `
			// 		resource "mezmo_pipeline" "test_parent" {
			// 			title = "parent pipeline"
			// 		}
			// 		resource "mezmo_demo_source" "my_source" {
			// 			pipeline = mezmo_pipeline.test_parent.id
			// 			format = "apache_common"
			// 		}`,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("mezmo_demo_source.my_source", "format", "apache_common"),
			// 		resource.TestCheckResourceAttr("mezmo_demo_source.my_source", "generation_id", "1"),
			// 	),
			// },
			// Delete testing automatically occurs in TestCase
		},
	})
}
