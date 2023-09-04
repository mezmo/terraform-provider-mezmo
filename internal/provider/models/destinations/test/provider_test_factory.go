package destinations

import (
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/mezmo-inc/terraform-provider-mezmo/internal/provider"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"mezmo": providerserver.NewProtocol6WithError(provider.New("destinations_test")()),
}

func init() {
	os.Setenv("TF_ACC", "1")
}
