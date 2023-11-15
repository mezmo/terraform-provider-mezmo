package sources

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestSQSSource(t *testing.T) {
	cacheKey := "sqs_source_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Error: Required field "pipeline_id"
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
							title = "parent pipeline"
						}`) + `
					resource "mezmo_sqs_source" "my_source" {
						region = "us-east-2"
						auth = {
							access_key_id = "123"
							secret_access_key = "secret123"
						}
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},
			// Error: Required field "queue_url"
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_sqs_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						region = "us-east-2"
						auth = {
							access_key_id = "123"
							secret_access_key = "secret123"
						}
					}`,
				ExpectError: regexp.MustCompile("The argument \"queue_url\" is required"),
			},
			// Error: Required field "auth"
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_sqs_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						queue_url = "http://example.com/queue"
						region = "us-east-2"
					}`,
				ExpectError: regexp.MustCompile("The argument \"auth\" is required"),
			},
			// Error: Required field "auth.access_key_id"
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_sqs_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						queue_url = "http://example.com/queue"
						region = "us-east-2"
						auth = {
							secret_access_key = "secret123"
						}
					}`,
				ExpectError: regexp.MustCompile("attribute \"access_key_id\" is\n\\s*required"),
			},
			// Error: Required field "auth.secret_access_key"
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_sqs_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						queue_url = "http://example.com/queue"
						region = "us-east-2"
						auth = {
							access_key_id = "123"
						}
					}`,
				ExpectError: regexp.MustCompile("attribute \"secret_access_key\" is\n\\s*required"),
			},
			// Error: Required field "region"
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_sqs_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						queue_url = "http://example.com/queue"
						auth = {
							access_key_id = "123"
							secret_access_key = "secret123"
						}
					}`,
				ExpectError: regexp.MustCompile("The argument \"region\" is required"),
			},

			// Create and Read testing
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_sqs_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						description = "my description"
						title = "my title"
						queue_url = "http://example.com/queue"
						region = "us-east-2"
						auth = {
							access_key_id = "123"
							secret_access_key = "secret123"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("mezmo_sqs_source.my_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_sqs_source.my_source", map[string]any{
						"description":            "my description",
						"title":                  "my title",
						"region":                 "us-east-2",
						"queue_url":              "http://example.com/queue",
						"auth.access_key_id":     "123",
						"auth.secret_access_key": "secret123",
						"generation_id":          "0",
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_sqs_source" "import_target" {
						pipeline_id = mezmo_pipeline.test_parent.id
						description = "my description"
						title = "my title"
						queue_url = "http://example.com/queue"
						region = "us-east-2"
						auth = {
							access_key_id = "123"
							secret_access_key = "secret123"
						}
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_sqs_source.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_sqs_source.my_source"),
				ImportStateVerify: true,
			},

			// Updates
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_sqs_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						description = "new description"
						title = "new title"
						queue_url = "http://example.com/another/queue"
						region = "us-east-1"
						auth = {
							access_key_id = "456"
							secret_access_key = "abc123"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("mezmo_sqs_source.my_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_sqs_source.my_source", map[string]any{
						"description":            "new description",
						"title":                  "new title",
						"region":                 "us-east-1",
						"queue_url":              "http://example.com/another/queue",
						"auth.access_key_id":     "456",
						"auth.secret_access_key": "abc123",
						"generation_id":          "1",
					}),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
