package providertest

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const authAccountId = "tf_test_01"

var setupMutex sync.Mutex
var isAccountCreated bool
var testConfigCache = make(map[string]string, 0)

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
				"x-auth-user-email" = "info@mezmo.com"
			}
		}
		`, GetTestEndpoint(), authAccountId)
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

		// For debugging:
		// fmt.Printf("---------- attributes ------- %+v\n", attributes)

		for expectedKey, expectedVal := range expected {
			foundVal, state_has_key := attributes[expectedKey]
			if state_has_key {
				lookupVal, err := lookupValue(expectedVal.(string), s)
				if err != nil {
					return err
				}
				if foundVal != lookupVal {
					return fmt.Errorf("Expected values do not match for key \"%s\". Found value: %s, Expected value: %s", expectedKey, foundVal, lookupVal)
				}
			} else if expectedVal != nil { // Using a nil value means the key should not be present at all
				return fmt.Errorf("Expected key \"%s\" was not found in %s", expectedKey, resourceName)
			}
		}
		return nil
	}
}
