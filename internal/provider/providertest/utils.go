package providertest

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"text/template"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
)

const authAccountId = "tf_test_01"
const authUserEmail = "info@mezmo.com"

var setupMutex sync.Mutex
var isAccountCreated bool
var testConfigCache = make(map[string]string, 0)

// IDRegex expression for Pipeline IDs
var IDRegex = regexp.MustCompile(`[\w-]{36}`)

// ParsedAccConfig applies config to tpl for acceptance tests
func ParsedAccConfig(config any, tpl string) (string, error) {
	var buf bytes.Buffer
	tmpl, _ := template.New("test").Parse(tpl)
	err := tmpl.Execute(&buf, config)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func SetCachedConfig(key string, body string) string {
	toCache := GetProviderConfig() + body
	testConfigCache[key] = toCache
	return toCache
}
func GetCachedConfig(key string) string {
	if fromCache, ok := testConfigCache[key]; ok {
		return fromCache
	}
	panic(fmt.Sprintf("GetCachedConfig cannot find key %s", key))
}

func GetTestEndpoint() string {
	endpoint := os.Getenv("TEST_ENDPOINT")
	if endpoint == "" {
		// Use port exposed in docker compose service
		endpoint = "http://localhost:19095"
	}
	return endpoint
}

func GetProviderConfig() string {
	return fmt.Sprintf(`
		provider "mezmo" {
			endpoint = %q
			auth_key = ""
			headers  = {
				// Used for authenticating against the service directly
				"x-auth-account-id"  = %q
				"x-auth-user-email" = %q
			}
		}
		`, GetTestEndpoint(), authAccountId, authUserEmail)
}

func TestPreCheck(t *testing.T) {
	defer setupMutex.Unlock()
	setupMutex.Lock()

	if isAccountCreated {
		return
	}

	controlToken := os.Getenv("TEST_CONTROL_TOKEN")
	client := http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest(
		http.MethodPut,
		GetTestEndpoint()+"/internal/account",
		strings.NewReader(fmt.Sprintf(`{"log_analysis_id": %q}`, authAccountId)))
	req.Header.Add("x-control-token", controlToken)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error while creating the account: %s", err.Error())
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("Error while reading body: %s", err.Error())
	}

	if resp.StatusCode >= http.StatusOK && resp.StatusCode <= http.StatusNoContent {
		t.Log("Created account for testing", string(body))
		isAccountCreated = true
	} else {
		t.Fatalf("Unexpected response when creating the test account: %s %s", resp.Status, string(body))
	}
}

var lookupRegex = regexp.MustCompile("^#(.+)\\.(.+)$")

// Given a resource name, look up its properties in state and return the requested value.
// Lookup syntax is supported for keys and must match #resource.property format.
func lookupValue(key string, s *terraform.State) (string, error) {
	matches := lookupRegex.FindStringSubmatch(key)
	length := len(matches)

	if length == 0 {
		return key, nil
	} else if length != 3 {
		err := fmt.Errorf("lookup pattern is not the correct structure: %v", matches[1:])
		return "", err
	}

	resourceName := matches[1]
	propertyName := matches[2]
	attributes := s.RootModule().Resources[resourceName].Primary.Attributes
	value := attributes[propertyName]

	if value == "" {
		err := fmt.Errorf("lookup for key \"%s\" found a blank value in attributes: %v", key, attributes)
		return "", err
	}
	return value, nil
}

// Given a string map, compare the expected values for those in the state.
// Lookup syntax is supported for map values, e.g "#mezmo_stringify_processor.my_processor.id"
func StateHasExpectedValues(resourceName string, expected map[string]any) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource := s.RootModule().Resources[resourceName]
		if resource == nil {
			return fmt.Errorf("resource \"%s\" not found", resourceName)
		}
		attributes := resource.Primary.Attributes

		if os.Getenv("DEBUG_ATTRIBUTES") == "1" {
			modelutils.PrintJSON(fmt.Sprintf("------ %s STATE ATTRIBUTES ------", resourceName), attributes)
		}

		for expectedKey, expectedVal := range expected {
			foundVal, state_has_key := attributes[expectedKey]
			if state_has_key {
				switch expectedVal := expectedVal.(type) {
				case *regexp.Regexp:
					if !expectedVal.MatchString(foundVal) {
						return fmt.Errorf("expected value \"%s\" for key \"%s\" does not match pattern \"%s\"", foundVal, expectedKey, expectedVal.String())
					}
				default:
					lookupVal, err := lookupValue(expectedVal.(string), s)
					if err != nil {
						return err
					}
					if foundVal != lookupVal {
						return fmt.Errorf("Expected values do not match for key \"%s\". Found value: %s, Expected value: %s", expectedKey, foundVal, lookupVal)
					}
				}
			} else if expectedVal != nil { // Using a nil value means the key should not be present at all
				return fmt.Errorf("Expected key \"%s\" was not found in %s", expectedKey, resourceName)
			}
		}
		return nil
	}
}

func StateDoesNotHaveFields(resourceName string, fields []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource := s.RootModule().Resources[resourceName]
		if resource == nil {
			return fmt.Errorf("resource \"%s\" not found", resourceName)
		}
		attributes := resource.Primary.Attributes

		if os.Getenv("DEBUG_ATTRIBUTES") == "1" {
			modelutils.PrintJSON(fmt.Sprintf("------ %s STATE ATTRIBUTES ------", resourceName), attributes)
		}

		for _, field := range fields {
			val, found := attributes[field]

			if found {
				return fmt.Errorf("Expected key \"%s\" not to exist in %s. Found %s", field, resourceName, val)
			}
		}

		return nil
	}
}

func ComputeImportId(resourceName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		resource := state.RootModule().Resources[resourceName]
		if resource == nil {
			return "", fmt.Errorf("resource \"%s\" not found", resourceName)
		}
		attributes := resource.Primary.Attributes
		if pipelineId, ok := attributes["pipeline_id"]; ok {
			return fmt.Sprintf("%s/%s", pipelineId, resource.Primary.ID), nil
		}
		return "", fmt.Errorf("resource \"%s\" does not have an attribute \"pipeline_id\"", resourceName)
	}
}

func TestDeletePipelineManually(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return deleteResource(
			resourceName,
			getResourceID(s, resourceName),
			"",
		)
	}
}

func TestDeletePipelineNodeManually(pipelineResourceName string, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return deleteResource(
			resourceName,
			getResourceID(s, pipelineResourceName),
			getResourceID(s, resourceName),
		)
	}
}

func getResourceID(s *terraform.State, resourceName string) string {
	res, ok := s.RootModule().Resources[resourceName]
	if ok {
		return res.Primary.ID
	}
	return ""
}

func deleteResource(resourceName string, pipelineId string, resourceId string) error {
	// resource names from TF state are of the form <resource_type>.<name>
	// example: mezmo_pipeline.test
	resourceType := strings.Split(resourceName, ".")[0]
	endpoint := fmt.Sprintf("/v3/pipeline/%s", pipelineId)
	if resourceType == "mezmo_pipeline" {
		err := makeDeleteRequest(endpoint)
		return err
	}
	nodeType, err := pipelineNodeType(resourceType)
	if err != nil {
		return err
	}
	if resourceId == "" {
		return fmt.Errorf("resource ID is required for %s deletion", resourceName)
	}
	// convert to v3/pipeline/<node_type>/<node_id>
	endpoint = fmt.Sprintf("%s/%s/%s", endpoint, nodeType, resourceId)
	err = makeDeleteRequest(endpoint)
	return err
}

func pipelineNodeType(resourceType string) (string, error) {
	switch {
	case strings.HasSuffix(resourceType, "_destination"):
		return "sink", nil
	case strings.HasSuffix(resourceType, "_processor"):
		return "transform", nil
	case strings.HasSuffix(resourceType, "_source"):
		return "source", nil
	default:
		return "", fmt.Errorf("unknown resource type: %s", resourceType)
	}
}

func makeDeleteRequest(urlPath string) error {
	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(
		http.MethodDelete,
		GetTestEndpoint()+urlPath,
		nil,
	)
	if err != nil {
		return fmt.Errorf("unable to create http request %s. reason: %s", urlPath, err.Error())
	}
	req.Header.Add("x-account-id", authAccountId)
	req.Header.Add("x-auth-user-email", authUserEmail)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request %s failed. reason: %s", urlPath, err.Error())
	}
	if resp.StatusCode >= http.StatusOK && resp.StatusCode <= http.StatusNoContent {
		return nil
	}
	return fmt.Errorf("received status code %d for %s", resp.StatusCode, urlPath)
}

func CheckMultipleErrors(err_strings []string) resource.ErrorCheckFunc {
	return func(err error) error {
		// Testing multiple regex errors is not possible with `ExpectError`, so we use
		// this custom function. Because this option is a `TestCase` option,
		// this test only has 1 `TestStep`.
		error_bytes := []byte(err.Error())

		for _, err_string := range err_strings {
			found, _ := regexp.Match(err_string, error_bytes)
			if !found {
				return errors.New("The expected error was not found: " + err_string)
			}
		}
		return nil
	}
}
