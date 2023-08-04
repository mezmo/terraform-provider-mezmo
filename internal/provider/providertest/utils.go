package providertest

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

const authAccountId = "tf_test_01"

var setupMutex sync.Mutex
var isAccountCreated bool

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
			auth_key        = %q
			endpoint        = %q
			auth_header     = "x-auth-account-id" // Used for authenticating against the service directly
			auth_additional = "info@mezmo.com"
		}
		`, authAccountId, GetTestEndpoint())
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
