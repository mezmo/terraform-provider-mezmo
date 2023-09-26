package destinations

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

func TestKafkaDestinationResource(t *testing.T) {
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
					resource "mezmo_kafka_destination" "my_destination" {
						itle = "my kafka title"
						description = "my kafka description"
					}`,
				ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
			},
			// Error: Missing brokers
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						topic = "topic1"
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
					resource "mezmo_kafka_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    port = 9092
						}]
						topic = "topic1"
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
					resource "mezmo_kafka_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						}]
						topic = "topic1"
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
					resource "mezmo_kafka_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 65536
						}]
						topic = "topic1"
						sasl = {
                            username = "my_username"
                            password = "my_password"
                            mechanism = "PLAIN"
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute brokers\\[0\\].port value must be between 1 and 65535, got: 65536"),
			},
			// Error: Missing topic
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						}]
						sasl = {
                            username = "my_username"
                            password = "my_password"
                            mechanism = "PLAIN"
						}
					}`,
				ExpectError: regexp.MustCompile("The argument \"topic\" is required, but no definition was found."),
			},
			// Error: Topic contains invalid topic (empty string)
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						}]
						topic = ""
						sasl = {
                            username = "my_username"
                            password = "my_password"
                            mechanism = "PLAIN"
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute topic string length must be at least 1, got: 0"),
			},
			// Error: Invalid event_key_id (empty string)
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						}]
						topic = "topic1"
						event_key_field = ""
						sasl = {
                            username = "my_username"
                            password = "my_password"
                            mechanism = "PLAIN"
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute event_key_field string length must be at least 1, got: 0"),
			},
			// Error: Invalid compression (non-matching enum value)
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						}]
						topic = "topic1"
						compression = "wrong"
						tls_enabled = true
						sasl = {
                            username = "my_username"
                            password = "my_password"
                            mechanism = "PLAIN"
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute compression value must be one of"),
			},
			// Error: Invalid encoding (non-matching enum value)
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						}]
						topic = "topic1"
						encoding = "wrong"
						sasl = {
                            username = "my_username"
                            password = "my_password"
                            mechanism = "PLAIN"
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute encoding value must be one of"),
			},
			// Error: Missing sasl username
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						}]
						topic = "topic1"
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
					resource "mezmo_kafka_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						}]
						topic = "topic1"
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
					resource "mezmo_kafka_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						}]
						topic = "topic1"
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
					resource "mezmo_kafka_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						}]
						topic = "topic1"
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
					resource "mezmo_kafka_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						}]
						topic = "topic1"
						sasl = {
                            username = "my_username"
                            password = "my_password"
                            mechanism = "wrong"
						}
					}`,
				ExpectError: regexp.MustCompile("Attribute sasl.mechanism value must be one of"),
			},
			// Create and Read testing
			{
				Config: SetCachedConfig(cacheKey, `
					resource "mezmo_pipeline" "test_parent" {
						title = "pipeline"
					}
					resource "mezmo_http_source" "my_source" {
						pipeline_id = mezmo_pipeline.test_parent.id
					}`) + `
					resource "mezmo_kafka_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						}]
						event_key_field = "my_key"
						topic = "topic1"
						compression = "gzip"
						encoding = "json"
						ack_enabled = true
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_kafka_destination.my_destination", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_kafka_destination.my_destination", map[string]any{
						"pipeline_id":     "#mezmo_pipeline.test_parent.id",
						"title":           "my kafka title",
						"description":     "my kafka description",
						"generation_id":   "0",
						"brokers.#":       "1",
						"brokers.0.host":  "mezmo.com",
						"brokers.0.port":  "9092",
						"event_key_field": "my_key",
						"topic":           "topic1",
						"tls_enabled":     "true",
						"compression":     "gzip",
						"encoding":        "json",
						"inputs.#":        "0",
						"ack_enabled":     "true",
					}),
				),
			},
			// Update and Read with new broker testing
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						},{
                            host = "mezmo.com"
                            port = 9093
						}]
						inputs = [mezmo_http_source.my_source.id]
						event_key_field = "my_key"
						topic = "topic1"
						compression = "gzip"
						encoding = "json"
						ack_enabled = true
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_kafka_destination.my_destination", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_kafka_destination.my_destination", map[string]any{
						"pipeline_id":     "#mezmo_pipeline.test_parent.id",
						"title":           "my kafka title",
						"description":     "my kafka description",
						"generation_id":   "1",
						"brokers.#":       "2",
						"brokers.0.host":  "mezmo.com",
						"brokers.0.port":  "9092",
						"brokers.1.host":  "mezmo.com",
						"brokers.1.port":  "9093",
						"inputs.#":        "1",
						"inputs.0":        "#mezmo_http_source.my_source.id",
						"event_key_field": "my_key",
						"topic":           "topic1",
						"tls_enabled":     "true",
						"compression":     "gzip",
						"encoding":        "json",
						"ack_enabled":     "true",
					}),
				),
			},
			// Update and Read with SASL testing
			{
				Config: GetCachedConfig(cacheKey) + `
					resource "mezmo_kafka_destination" "my_destination" {
						pipeline_id = mezmo_pipeline.test_parent.id
						title = "my kafka title"
						description = "my kafka description"
						brokers = [{
						    host = "mezmo.com"
						    port = 9092
						},{
                            host = "mezmo.com"
                            port = 9093
						}]
						inputs = [mezmo_http_source.my_source.id]
						event_key_field = "my_key"
						topic = "topic1"
						compression = "gzip"
						encoding = "json"
						ack_enabled = true
						sasl = {
                            username = "my_username"
                            password = "my_password"
                            mechanism = "PLAIN"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"mezmo_kafka_destination.my_destination", "id", regexp.MustCompile(`[\w-]{36}`)),
					StateHasExpectedValues("mezmo_kafka_destination.my_destination", map[string]any{
						"pipeline_id":     "#mezmo_pipeline.test_parent.id",
						"title":           "my kafka title",
						"description":     "my kafka description",
						"generation_id":   "2",
						"brokers.#":       "2",
						"brokers.0.host":  "mezmo.com",
						"brokers.0.port":  "9092",
						"brokers.1.host":  "mezmo.com",
						"brokers.1.port":  "9093",
						"event_key_field": "my_key",
						"topic":           "topic1",
						"tls_enabled":     "true",
						"sasl.username":   "my_username",
						"sasl.password":   "my_password",
						"sasl.mechanism":  "PLAIN",
						"compression":     "gzip",
						"encoding":        "json",
						"ack_enabled":     "true",
					}),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
