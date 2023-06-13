package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPipelineResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: getProviderConfig() + `
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
			// {
			// 	Config: getProviderConfig() + `
			// 		resource "mezmo_pipeline" "test" {
			// 			title = "updated title"
			// 		}`,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("mezmo_pipeline.test", "title", "updated title"),
			// 	),
			// },
			// Delete testing automatically occurs in TestCase
		},
	})
}
