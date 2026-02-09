package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/v5/internal/provider/providertest"
)

// Note that this test is not able to test all the possible use cases for this resource.
// There currently isn't any information on how to unit test with submodules.
// Having other resources in line with `publish_pipeline` will cause race conditions and unintended behavior.
func TestPublishPipelineResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Cache base resources of pipeline and source
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "my_pipeline" {
						title = "pipeline"
					}
					resource "mezmo_publish_pipeline" "my_publish_pipeline" {
						pipeline_id = mezmo_pipeline.my_pipeline.id
          }
					`,
				ExpectNonEmptyPlan: true, // We always re-create which causes a non-empty plan.
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mezmo_pipeline.my_pipeline", "title", "pipeline"),
					StateHasExpectedValues("mezmo_publish_pipeline.my_publish_pipeline", map[string]any{
						"pipeline_id": "#mezmo_pipeline.my_pipeline.id",
					}),
				),
			},
			// Update
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "my_pipeline" {
						title = "Updated Pipeline"
					}
					resource "mezmo_publish_pipeline" "my_publish_pipeline" {
						pipeline_id = mezmo_pipeline.my_pipeline.id
          }
					`,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mezmo_pipeline.my_pipeline", "title", "Updated Pipeline"),
					StateHasExpectedValues("mezmo_publish_pipeline.my_publish_pipeline", map[string]any{
						"pipeline_id": "#mezmo_pipeline.my_pipeline.id",
					}),
				),
			},
		},
	})
}
