package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/providertest"
)

func TestPipelineResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Required fields test
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test" {
					}`,
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			// Create and Read testing
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test" {
						title = "hello title"
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify user-defined properties
					resource.TestCheckResourceAttr("mezmo_pipeline.test", "title", "hello title"),
					// Verify computed properties
					resource.TestCheckResourceAttrSet("mezmo_pipeline.test", "id"),
					resource.TestCheckResourceAttrSet("mezmo_pipeline.test", "created_at"),
				),
			},
			// Update and Read testing
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test" {
						title = "updated title"
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mezmo_pipeline.test", "title", "updated title"),
				),
			},
			// manually delete a pipeline and verify it is re-created
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_deleted" {
						title = "updated title"
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mezmo_pipeline.test_deleted", "title", "updated title"),
					// delete the resource
					TestDeletePipelineManually("mezmo_pipeline.test_deleted"),
				),
				// verify resource is created after refresh
				ExpectNonEmptyPlan: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
