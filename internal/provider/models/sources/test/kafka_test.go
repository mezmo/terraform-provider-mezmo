package sources

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/providertest"
)

func TestKafkaSourceResource(t *testing.T) {
	const cacheKey = "kafka_resources"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestPreCheck(t) },
		Steps: []resource.TestStep{
			// Error: pipeline_id is required
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}`) + `
					resource "mezmo_kafka_source" "my_source" {
						title = "my kafka title"
						description = "my kafka description"
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},
			// Error: Missing brokers
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						topics = ["topic1", "topic2"]
						group_id = "my_group_id"
						tls_enabled = true
						sasl = {
                            username = "my_username"
                            password = "my_password"
                            mechanism = "PLAIN"
						}
					}`,
				ExpectError: regexp.MustCompile("The argument \"brokers\" is required, but no definition was found."),
			},
			// Error: Broker missing host
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    port = 9092
						}]
						topics = ["topic1", "topic2"]
						group_id = "my_group_id"
						tls_enabled = true
						sasl = {
                            username = "my_username"
                            password = "my_password"
                            mechanism = "PLAIN"
						}
					}`,
				ExpectError: regexp.MustCompile("element 0: attribute \"host\" is\\s*required."),
			},
			// Error: Broker missing port
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						}]
						topics = ["topic1", "topic2"]
						group_id = "my_group_id"
						tls_enabled = true
						sasl = {
                            username = "my_username"
                            password = "my_password"
                            mechanism = "PLAIN"
						}
					}`,
				ExpectError: regexp.MustCompile("element 0: attribute \"port\" is\\s*required."),
			},
			// Error: Invalid broker port (> 65535)
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 65536
						}]
						topics = ["topic1", "topic2"]
						group_id = "my_group_id"
						tls_enabled = true
						sasl = {
                            username = "my_username"
                            password = "my_password"
                            mechanism = "PLAIN"
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute brokers\\[0\\].port value must be between 1 and 65535, got: 65536"),
			},
			// Error: Missing topics
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						}]
						group_id = "my_group_id"
						tls_enabled = true
						sasl = {
                            username = "my_username"
                            password = "my_password"
                            mechanism = "PLAIN"
						}
					}`,
				ExpectError: regexp.MustCompile("The argument \"topics\" is required, but no definition was found."),
			},
			// Error: Topics contains invalid topic (empty string)
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						}]
						topics = [""]
						group_id = "my_group_id"
						tls_enabled = true
						sasl = {
                            username = "my_username"
                            password = "my_password"
                            mechanism = "PLAIN"
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute topics\\[0\\] string length must be between 1 and 256, got: 0"),
			},
			// Error: Topics contains invalid topic (too long string)
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						}]
						topics = ["aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"]
						group_id = "my_group_id"
						tls_enabled = true
						sasl = {
                            username = "my_username"
                            password = "my_password"
                            mechanism = "PLAIN"
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute topics\\[0\\] string length must be between 1 and 256, got: 260"),
			},
			// Error: Missing group_id
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						}]
						tls_enabled = true
						sasl = {
                            username = "my_username"
                            password = "my_password"
                            mechanism = "PLAIN"
						}
					}`,
				ExpectError: regexp.MustCompile("The argument \"group_id\" is required, but no definition was found."),
			},
			// Error: Invalid group_id (empty string)
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						}]
						topics = ["topic1", "topic2"]
						group_id = ""
						tls_enabled = true
						sasl = {
                            username = "my_username"
                            password = "my_password"
                            mechanism = "PLAIN"
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute group_id string length must be at least 1, got: 0"),
			},
			// Error: Missing sasl username
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						}]
						topics = ["topic1", "topic2"]
						group_id = "my_group_id"
						tls_enabled = true
						sasl = {
                            password = "my_password"
                            mechanism = "PLAIN"
						}
					}`,
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"sasl\": attribute \"username\" is required."),
			},
			// Error: Invalid sasl username (empty string)
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						}]
						topics = ["topic1", "topic2"]
						group_id = "my_group_id"
						tls_enabled = true
						sasl = {
						    username = ""
                            password = "my_password"
                            mechanism = "PLAIN"
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute sasl.username string length must be at least 1, got: 0"),
			},
			// Error: Missing sasl password
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						}]
						topics = ["topic1", "topic2"]
						group_id = "my_group_id"
						tls_enabled = true
						sasl = {
						    username = "my_username"
                            mechanism = "PLAIN"
						}
					}`,
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"sasl\": attribute \"password\" is required."),
			},
			// Error: Invalid sasl password (empty string)
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						}]
						topics = ["topic1", "topic2"]
						group_id = "my_group_id"
						tls_enabled = true
						sasl = {
						    username = "my_username"
                            password = ""
                            mechanism = "PLAIN"
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute sasl.password string length must be at least 1, got: 0"),
			},
			// Error: invalid sasl mechanism (non-matching enum)
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						}]
						topics = ["topic1", "topic2"]
						group_id = "my_group_id"
						tls_enabled = true
						sasl = {
                            username = "my_username"
                            password = "my_password"
                            mechanism = "wrong"
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute sasl.mechanism value must be one of"),
			},
			// Error: Invalid decoding (non-matching enum)
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						}]
						topics = ["topic1", "topic2"]
						group_id = "my_group_id"
						tls_enabled = true
						sasl = {
						    username = "my_username"
                            password = "my_password"
                            mechanism = "PLAIN"
						}
						decoding = "wrong"
					}`,
				ExpectError: regexp.MustCompile("Attribute decoding value must be one of"),
			},
			// Create and Read testing
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						}]
						topics = ["topic1", "topic2"]
						group_id = "my_group_id"
						tls_enabled = true
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_kafka_source.my_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_kafka_source.my_source", map[string]any{
						"pipeline_id":    "#mezmo_pipeline.test_parent.id",
						"title":          "my kafka title",
						"description":    "my kafka description",
						"generation_id":  "0",
						"brokers.#":      "1",
						"brokers.0.host": "mezmo.com",
						"brokers.0.port": "9092",
						"topics.#":       "2",
						"topics.0":       "topic1",
						"topics.1":       "topic2",
						"group_id":       "my_group_id",
						"tls_enabled":    "true",
					}),
				),
			},
			// Update and Read with new broker testing
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						},
						{
                            host = "mezmo.com"
                            port = 9093
						}]
						topics = ["topic1", "topic2"]
						group_id = "my_group_id"
						tls_enabled = true
						sasl = {
                            username = "my_username"
                            password = "my_password"
                            mechanism = "PLAIN"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_kafka_source.my_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_kafka_source.my_source", map[string]any{
						"pipeline_id":    "#mezmo_pipeline.test_parent.id",
						"title":          "my kafka title",
						"description":    "my kafka description",
						"generation_id":  "1",
						"brokers.#":      "2",
						"brokers.0.host": "mezmo.com",
						"brokers.0.port": "9092",
						"brokers.1.host": "mezmo.com",
						"brokers.1.port": "9093",
						"topics.#":       "2",
						"topics.0":       "topic1",
						"topics.1":       "topic2",
						"group_id":       "my_group_id",
						"tls_enabled":    "true",
						"sasl.username":  "my_username",
						"sasl.password":  "my_password",
						"sasl.mechanism": "PLAIN",
					}),
				),
			},
			// Update and Read with SASL testing
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						}]
						topics = ["topic1", "topic2"]
						group_id = "my_group_id"
						tls_enabled = true
						sasl = {
                            username = "my_username"
                            password = "my_password"
                            mechanism = "PLAIN"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_kafka_source.my_source", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_kafka_source.my_source", map[string]any{
						"pipeline_id":    "#mezmo_pipeline.test_parent.id",
						"title":          "my kafka title",
						"description":    "my kafka description",
						"generation_id":  "2",
						"brokers.#":      "1",
						"brokers.0.host": "mezmo.com",
						"brokers.0.port": "9092",
						"topics.#":       "2",
						"topics.0":       "topic1",
						"topics.1":       "topic2",
						"group_id":       "my_group_id",
						"tls_enabled":    "true",
						"sasl.username":  "my_username",
						"sasl.password":  "my_password",
						"sasl.mechanism": "PLAIN",
					}),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
