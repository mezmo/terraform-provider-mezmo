package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mezmo/terraform-provider-mezmo/v4/internal/client"
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
					func(s *terraform.State) error {

						// Read mezmo_pipeline.test local state
						model, ok := s.RootModule().Resources["mezmo_pipeline.test"]

						if !ok {
							return fmt.Errorf("PipelineResourceModel not found in terraform.State")
						}

						// Get id and ctx to call Mezmo API in test with the NewTestClient constructor
						id := model.Primary.Attributes["id"]
						ctx := context.Background()
						c := NewTestClient()

						// Pipeline from Mezmo API call
						pipeline, api_err := c.Pipeline(id, ctx)

						if api_err != nil {
							return fmt.Errorf("Pipeline id %s not found in mezmo api call", id)
						}

						// Verify origin in remote Mezmo API is indeed set to "terraform"
						if pipeline.Origin != client.ORIGIN_TERRAFORM {
							return fmt.Errorf("Pipeline created with origin set to \"%v\" but expected \"terraform\"", pipeline.Origin)
						}

						return nil
					},
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
