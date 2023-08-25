package sinks

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/providertest"
)

func TestNewRelicSinkResource(t *testing.T) {
	const cacheKey = "new_relic_sink_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Error: properties are required
			{
				Config: GetProviderConfig() + `
					resource "mezmo_new_relic_sink" "my_sink" {
						inputs      = ["abc"]
						account_id = "acc1"
						license_key = "key1"
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_new_relic_sink" "my_sink" {
						inputs      = ["abc"]
						pipeline_id = "pip1"
						license_key = "key1"
					}`,
				ExpectError: regexp.MustCompile("The argument \"account_id\" is required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_new_relic_sink" "my_sink" {
						inputs      = ["abc"]
						pipeline_id = "pip1"
						account_id = "acc1"
					}`,
				ExpectError: regexp.MustCompile("The argument \"license_key\" is required"),
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
					resource "mezmo_new_relic_sink" "my_sink" {
						title = "my sink"
						description = "my sink description"
						inputs      = [mezmo_http_source.my_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id
						account_id = "acc1"
						license_key = "key1"
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_new_relic_sink.my_sink", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_new_relic_sink.my_sink", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "my sink",
						"description":   "my sink description",
						"generation_id": "0",
						"ack_enabled":   "true",
						"inputs.#":      "1",
						"account_id":    "acc1",
						"license_key":   "key1",
						"api":           "logs",
					}),
				),
			},

			// Update all fields (pass through)
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_new_relic_sink" "my_sink" {
						title = "new title"
						description = "new description"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs      = [mezmo_http_source.my_source.id]
						account_id  = "acc2"
						license_key = "key2"
						api         = "metrics"
						ack_enabled = false
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_new_relic_sink.my_sink", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "new title",
						"description":   "new description",
						"generation_id": "1",
						"inputs.#":      "1",
						"inputs.0":      "#mezmo_http_source.my_source.id",
						"ack_enabled":   "false",
						"account_id":    "acc2",
						"license_key":   "key2",
						"api":           "metrics",
					}),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
