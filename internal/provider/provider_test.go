package provider

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const authAccountId = "tf_test_01"
const endpoint = "http://localhost:19095"

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"mezmo": providerserver.NewProtocol6WithError(New("test")()),
}

var setupMutex sync.Mutex
var isAccountCreated bool

// Terraform uses a flag to prevent tests from running unintentionally,
// as Terraform resources are often linked to real-world resources and infrastructure.
// In our case, we support running the integration tests against a pipeline-service container.
// We could support running against dev / staging in the future, once we expose the pipeline-service
// behind the Gateway.
func init() {
	os.Setenv("TF_ACC", "1")
}

func getProviderConfig() string {
	return fmt.Sprintf(`
		provider "mezmo" {
			auth_key = %q
			endpoint = %q
			auth_header = "x-auth-account-id" // Use for authenticating against the service directly
		}
		`, authAccountId, endpoint)
}

func testAccPreCheck(t *testing.T) {
	defer setupMutex.Unlock()
	setupMutex.Lock()

	if isAccountCreated {
		return
	}

	controlToken := os.Getenv("CONTROL_TOKEN")
	client := http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest(
		http.MethodPut,
		endpoint+"/v1/control/account",
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
