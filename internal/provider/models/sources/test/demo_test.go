package sources

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mezmo-inc/terraform-provider-mezmo/internal/provider"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/providertest"
)

func init() {
	os.Setenv("TF_ACC", "1")
}

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"mezmo": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func TestDemoSourceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Required fields json
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_demo_source" "my_source" {
						pipeline = mezmo_pipeline.test_parent.id
					}`,
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			// Required fields parent pipeline id
			{
				Config: GetProviderConfig() + `
					resource "mezmo_demo_source" "my_source" {
						format = "json"
					}`,
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			// Required fields parent pipeline id
			{
				Config: GetProviderConfig() + `
					resource "mezmo_demo_source" "my_source" {
						pipeline = "798e1028-0b60-11ee-be56-0242ac120002"
						format = "NOT_VALID"
					}`,
				ExpectError: regexp.MustCompile("Attribute format value must be one of"),
			},
			// Create and Read testing
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_demo_source" "my_source" {
						pipeline = mezmo_pipeline.test_parent.id
						format = "json"
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify user-defined properties
					resource.TestCheckResourceAttr("mezmo_demo_source.my_source", "format", "json"),
					// Verify computed properties
					resource.TestCheckResourceAttrSet("mezmo_demo_source.my_source", "id"),
					resource.TestCheckResourceAttrSet("mezmo_demo_source.my_source", "generation_id"),
				),
			},
			// Update and Read testing
			{
				Config: GetProviderConfig() + `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_demo_source" "my_source" {
						pipeline = mezmo_pipeline.test_parent.id
						format = "apache_common"
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mezmo_demo_source.my_source", "format", "apache_common"),
					resource.TestCheckResourceAttr("mezmo_demo_source.my_source", "generation_id", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
