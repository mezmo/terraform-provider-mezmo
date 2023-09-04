package destinations

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/providertest"
)

func TestSplunkHecLogsDestinationResource(t *testing.T) {
	const cacheKey = "splunk_hec_logs_destination_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Error: properties are required
			{
				Config: GetProviderConfig() + `
					resource "mezmo_splunk_hec_logs_destination" "my_destination" {
						inputs   = ["abc"]
						endpoint = "http://example.com"
					}`,
				ExpectError: regexp.MustCompile("The argument \"token\" is required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_splunk_hec_logs_destination" "my_destination" {
						inputs = ["abc"]
						token  = "my_token"
					}`,
				ExpectError: regexp.MustCompile("The argument \"endpoint\" is required"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_splunk_hec_logs_destination" "my_destination" {
						inputs      = ["abc"]
						pipeline_id = "2ee4d436-466a-11ee-be56-0242ac120002"
						endpoint    = "http://example.com"
						token       = "my_token"
						source      = {}
					}`,
				ExpectError: regexp.MustCompile("source requires field or value to be defined"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_splunk_hec_logs_destination" "my_destination" {
						inputs      = ["abc"]
						pipeline_id = "2ee4d436-466a-11ee-be56-0242ac120002"
						endpoint    = "http://example.com"
						token       = "my_token"
						source_type = {}
					}`,
				ExpectError: regexp.MustCompile("source_type requires field or value to be defined"),
			},
			{
				Config: GetProviderConfig() + `
					resource "mezmo_splunk_hec_logs_destination" "my_destination" {
						inputs      = ["abc"]
						pipeline_id = "2ee4d436-466a-11ee-be56-0242ac120002"
						endpoint    = "http://example.com"
						token       = "my_token"
						index       = {}
					}`,
				ExpectError: regexp.MustCompile("index requires field or value to be defined"),
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
					resource "mezmo_splunk_hec_logs_destination" "my_destination" {
						title = "My destination"
						description = "my destination description"
						pipeline_id = mezmo_pipeline.test_parent.id
						endpoint    = "https://example.com"
						token       = "my_token"
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_splunk_hec_logs_destination.my_destination", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_splunk_hec_logs_destination.my_destination", map[string]any{
						"pipeline_id":            "#mezmo_pipeline.test_parent.id",
						"title":                  "My destination",
						"description":            "my destination description",
						"generation_id":          "0",
						"ack_enabled":            "true",
						"inputs.#":               "0",
						"endpoint":               "https://example.com",
						"token":                  "my_token",
						"compression":            "none",
						"tls_verify_certificate": "true",
						"host_field":             "metadata.host",
						"timestamp_field":        "metadata.time",
					}),
				),
			},

			// Update source fields and others
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_splunk_hec_logs_destination" "my_destination" {
						title = "my new destination"
						description = "my new destination description"
						inputs      = [mezmo_http_source.my_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id
						endpoint    = "https://example2.com"
						token       = "my_token2"
						ack_enabled = false
						tls_verify_certificate = false
						source      = {
							field = ".source"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_splunk_hec_logs_destination.my_destination", map[string]any{
						"pipeline_id":            "#mezmo_pipeline.test_parent.id",
						"title":                  "my new destination",
						"description":            "my new destination description",
						"generation_id":          "1",
						"ack_enabled":            "false",
						"inputs.#":               "1",
						"endpoint":               "https://example2.com",
						"token":                  "my_token2",
						"compression":            "none",
						"tls_verify_certificate": "false",
						"host_field":             "metadata.host",
						"timestamp_field":        "metadata.time",
						"source.field":           ".source",
					}),
				),
			},

			// Update other fields and nullify source
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_splunk_hec_logs_destination" "my_destination" {
						inputs      = [mezmo_http_source.my_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id
						endpoint    = "https://example3.com"
						token       = "my_token3"
						source_type = {
							field = ".my_source_type"
						}
						index = {
							field = ".my_index"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_splunk_hec_logs_destination.my_destination", map[string]any{
						"pipeline_id":            "#mezmo_pipeline.test_parent.id",
						"generation_id":          "2",
						"ack_enabled":            "true",
						"inputs.#":               "1",
						"endpoint":               "https://example3.com",
						"token":                  "my_token3",
						"compression":            "none",
						"tls_verify_certificate": "true",
						"host_field":             "metadata.host",
						"timestamp_field":        "metadata.time",
						"source.%":               nil,
						"source_type.field":      ".my_source_type",
						"source_type.value":      nil,
						"index.field":            ".my_index",
						"index.value":            nil,
					}),
				),
			},

			// Update set all values
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_splunk_hec_logs_destination" "my_destination" {
						inputs      = [mezmo_http_source.my_source.id]
						pipeline_id = mezmo_pipeline.test_parent.id
						endpoint    = "https://example4.com"
						token       = "my_token4"
						source = {
							value = "my fixed source"
						}
						source_type = {
							value = "my fixed source type"
						}
						index = {
							value = "my fixed index"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_splunk_hec_logs_destination.my_destination", map[string]any{
						"pipeline_id":            "#mezmo_pipeline.test_parent.id",
						"generation_id":          "3",
						"ack_enabled":            "true",
						"inputs.#":               "1",
						"endpoint":               "https://example4.com",
						"token":                  "my_token4",
						"compression":            "none",
						"tls_verify_certificate": "true",
						"host_field":             "metadata.host",
						"timestamp_field":        "metadata.time",
						"source.field":           nil,
						"source.value":           "my fixed source",
						"source_type.field":      nil,
						"source_type.value":      "my fixed source type",
						"index.field":            nil,
						"index.value":            "my fixed index",
					}),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
