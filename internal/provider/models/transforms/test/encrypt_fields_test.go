package transforms

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/providertest"
)

func TestEncryptFieldsTransform(t *testing.T) {
	const cacheKey = "encrypt_fields_resources"
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
					resource "mezmo_encrypt_fields_transform" "my_transform" {
						fields = [".nope"]
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},

			// Error: `field` is required
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_encrypt_fields_transform" "my_transform" {
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
					resource "mezmo_encrypt_fields_transform" "my_transform" {
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
					resource "mezmo_encrypt_fields_transform" "my_transform" {
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
					resource "mezmo_encrypt_fields_transform" "my_transform" {
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
					resource "mezmo_encrypt_fields_transform" "my_transform" {
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
				resource "mezmo_encrypt_fields_transform" "my_transform" {
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
				resource "mezmo_encrypt_fields_transform" "my_transform" {
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
				resource "mezmo_encrypt_fields_transform" "my_transform" {
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
					resource "mezmo_encrypt_fields_transform" "my_transform" {
						title = "encrypt fields title"
						description = "encrypt fields desc"
						pipeline_id = mezmo_pipeline.test_parent.id
						algorithm = "AES-128-CFB"
						key = "1111111111111111"
						iv_field = ".some_iv_field"
						field = ".something"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_encrypt_fields_transform.my_transform", "id", regexp.MustCompile(`[\w-]{36}`)),

					StateHasExpectedValues("mezmo_encrypt_fields_transform.my_transform", map[string]any{
						"pipeline_id":   "#mezmo_pipeline.test_parent.id",
						"title":         "encrypt fields title",
						"description":   "encrypt fields desc",
						"generation_id": "0",
						"inputs.#":      "0",
						"algorithm":     "AES-128-CFB",
						"key":           "1111111111111111",
						"iv_field":      ".some_iv_field",
						"field":         ".something",
					}),
				),
			},

			// Update fields
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_encrypt_fields_transform" "my_transform" {
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
					StateHasExpectedValues("mezmo_encrypt_fields_transform.my_transform", map[string]any{
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
				resource "mezmo_encrypt_fields_transform" "my_transform" {
					pipeline_id = mezmo_pipeline.test_parent.id
					inputs = [mezmo_http_source.my_source.id]
					algorithm = "AES-256-CBC-PKCS7"
					key = "key not long enough"
					iv_field = ".other_iv_field"
					field = ".something_else"
				}`,
				ExpectError: regexp.MustCompile("NOT have fewer than 32"),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
