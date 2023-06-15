package provider

import (
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"mezmo": providerserver.NewProtocol6WithError(New("test")()),
}

// Terraform uses a flag to prevent tests from running unintentionally,
// as Terraform resources are often linked to real-world resources and infrastructure.
// In our case, we support running the integration tests against a pipeline-service container.
// We could support running against dev / staging in the future, once we expose the pipeline-service
// behind the Gateway.
func init() {
	os.Setenv("TF_ACC", "1")
}
