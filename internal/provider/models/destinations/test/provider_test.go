package destinations

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mezmo/terraform-provider-mezmo/v4/internal/provider"
	"github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/providertest"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"mezmo": providerserver.NewProtocol6WithError(provider.New("destinations_test")()),
}

// testAccBackend checks that userConfig matches whats is in the
// backend for resourceName
func testAccBackend(resourceName string, userConfig map[string]any) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Minute)
		defer cancel()

		client := providertest.NewTestClient()
		destination, err := client.Destination(rs.Primary.Attributes["pipeline_id"], rs.Primary.ID, ctx)
		if err != nil {
			return err
		}

		// Check the backend UserConfig matches the passed userConfig
		return providertest.ValidateUserConfig(destination.UserConfig, userConfig)
	}
}
