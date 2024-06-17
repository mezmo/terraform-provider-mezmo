package processors

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestAccDecryptFieldsProcessor(t *testing.T) {
	const cacheKey = "decrypt_fields_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Error: pipeline_id is required
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_http_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`) + `
					resource "mezmo_decrypt_fields_processor" "my_processor" {
						fields = [".nope"]
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},

			// Error: `field` is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_decrypt_fields_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						algorithm = "AES-256-CFB"
						key = "mybiglongsecretatleast16chars"
						iv_field = ".some_iv_field"
					}`,
				ExpectError: regexp.MustCompile("The argument \"field\" is required"),
			},

			// Error: `algorithm` is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_decrypt_fields_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						key = "mybiglongsecretatleast16chars"
						iv_field = ".some_iv_field"
						field = ".something"
					}`,
				ExpectError: regexp.MustCompile("The argument \"algorithm\" is required"),
			},

			// Error: `algorithm` is an enum
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_decrypt_fields_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						key = "mybiglongsecretatleast16chars"
						iv_field = ".some_iv_field"
						field = ".something"
						algorithm = "Pig-Latin"
					}`,
				ExpectError: regexp.MustCompile("Attribute algorithm value must be one of:"),
			},

			// Error: `key` is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_decrypt_fields_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						algorithm = "AES-256-CFB"
						iv_field = ".some_iv_field"
						field = ".something"
					}`,
				ExpectError: regexp.MustCompile("The argument \"key\" is required"),
			},

			// Error: `iv_field` is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_decrypt_fields_processor" "my_processor" {
						pipeline_id = mezmo_pipeline.test_parent.id
						key = "mybiglongsecretatleast16chars"
						algorithm = "AES-256-CFB"
						field = ".something"
					}`,
				ExpectError: regexp.MustCompile("The argument \"iv_field\" is required"),
			},

			// Error: `field` values validates length
			{
				Config: GetCachedConfig(cacheKey) + `
				resource "mezmo_decrypt_fields_processor" "my_processor" {
					pipeline_id = mezmo_pipeline.test_parent.id
					algorithm = "AES-256-CFB"
					key = "mybiglongsecretatleast16chars"
					iv_field = ".some_iv_field"
					field = ""
				}`,
				ExpectError: regexp.MustCompile("Attribute field string length must be at least 1"),
			},

			// Error: `key` value too short
			{
				Config: GetCachedConfig(cacheKey) + `
				resource "mezmo_decrypt_fields_processor" "my_processor" {
					pipeline_id = mezmo_pipeline.test_parent.id
					algorithm = "AES-256-CFB"
					key = "not_enough"
					iv_field = ".some_iv_field"
					field = ".something"
				}`,
				ExpectError: regexp.MustCompile("Attribute key string length must be at least 16, got: 10"),
			},

			// Error: `key` value too long
			{
				Config: GetCachedConfig(cacheKey) + `
				resource "mezmo_decrypt_fields_processor" "my_processor" {
					pipeline_id = mezmo_pipeline.test_parent.id
					algorithm = "AES-256-CFB"
					key = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
					iv_field = ".some_iv_field"
					field = ".something"
				}`,
				ExpectError: regexp.MustCompile("Attribute key string length must be at most 32, got: 72"),
			},

			// Create with defaults
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_decrypt_fields_processor" "my_processor" {
						title = "decrypt fields title"
						description = "decrypt fields desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						algorithm = "AES-128-CFB"
						key = "1111111111111111"
						iv_field = ".some_iv_field"
						field = ".something"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_decrypt_fields_processor.my_processor", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_decrypt_fields_processor.my_processor", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "decrypt fields title",
						"description":   "decrypt fields desc",
						"generation_id": "0",
						"inputs.#":      "0",
						"algorithm":     "AES-128-CFB",
						"key":           "1111111111111111",
						"iv_field":      ".some_iv_field",
						"field":         ".something",
					}),
				),
			},

			// Import
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_decrypt_fields_processor" "import_target" {
						title = "decrypt fields title"
						description = "decrypt fields desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						algorithm = "AES-128-CFB"
						key = "1111111111111111"
						iv_field = ".some_iv_field"
						field = ".something"
					}`,
				ImportState:       true,
				ResourceName:      "mezmo_decrypt_fields_processor.import_target",
				ImportStateIdFunc: ComputeImportId("mezmo_decrypt_fields_processor.my_processor"),
				ImportStateVerify: true,
			},

			// Update fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_decrypt_fields_processor" "my_processor" {
						title = "new title"
						description = "new desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						inputs = [mezmo_http_source.my_source.id]
						algorithm = "AES-128-CTR"
						key = "2222222222222222"
						iv_field = ".other_iv_field"
						field = ".something_else"
					}`,
				Check: resource.ComposeTestCheckFunc(
					StateHasExpectedValues("mezmo_decrypt_fields_processor.my_processor", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "new title",
						"description":   "new desc",
						"generation_id": "1",
						"inputs.#":      "1",
						"inputs.0":      "#mezmo_http_source.my_source.id",
						"algorithm":     "AES-128-CTR",
						"key":           "2222222222222222",
						"iv_field":      ".other_iv_field",
						"field":         ".something_else",
					}),
				),
			},

			// Error: server-side validation
			{
				Config: GetCachedConfig(cacheKey) + `
				resource "mezmo_decrypt_fields_processor" "my_processor" {
					pipeline_id = mezmo_pipeline.test_parent.id
					inputs = [mezmo_http_source.my_source.id]
					algorithm = "AES-256-CBC-PKCS7"
					key = "key not long enough"
					iv_field = ".other_iv_field"
					field = ".something_else"
				}`,
				ExpectError: regexp.MustCompile("NOT have fewer than 32"),
			},
			// confirm manually deleted resources are recreated
			{
				Config: GetProviderConfig() + `
				resource "mezmo_pipeline" "test_parent2" {
					title = "pipeline"
				}
				resource "mezmo_decrypt_fields_processor" "test_processor" {
					pipeline_id = mezmo_pipeline.test_parent2.id
					title 			= "new title"
					inputs 			= []
					algorithm 	= "AES-128-CFB"
					key 				= "1111111111111111"
					iv_field 		= ".some_iv_field"
					field 			= ".something"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_decrypt_fields_processor.test_processor", "id", regexp.MustCompile(`[\w-]{36}`)),
					resource.TestCheckResourceAttr("mezmo_decrypt_fields_processor.test_processor", "title", "new title"),
					// delete the resource
					TestDeletePipelineNodeManually(
						"mezmo_pipeline.test_parent2",
						"mezmo_decrypt_fields_processor.test_processor",
					),
				),
				// verify resource will be re-created after refresh
				ExpectNonEmptyPlan: true,
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
